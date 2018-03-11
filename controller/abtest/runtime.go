//@author wuyong
//@date   2018/1/11
//@desc

package abtestController

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gateway-admin/lib"
	_ "gateway-admin/model"
	"gateway-admin/model/abtest/dao"
	"gateway-admin/model/abtest/serverpage"

	"github.com/gin-gonic/gin"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
)

func GetRuntimeList(context *gin.Context) {
	status := context.Query("status")
	isDelete := context.Query("is_delete")
	offset, _ := strconv.Atoi(context.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(context.DefaultQuery("limit", "100"))

	query := map[string]string{}
	if status != "" {
		query["status"] = status
	}
	if isDelete != "" {
		query["is_delete"] = isDelete
	}

	sortBy := []string{
		"create_time",
	}
	order := []string{
		"desc",
	}

	rows, err := abtestServerPage.GetRutimeInfo(query, sortBy, order, int64(offset), int64(limit))
	if err != nil {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}
	total, err := abtestDao.GetRuntimeTotalNum()
	if err != nil {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}

	if len(rows) > 0 {
		data := struct {
			TotalCount int64       `json:"totalCount"`
			List       interface{} `json:"list"`
		}{
			total,
			rows,
		}
		jsonData := lib.GetSuccess(data)
		context.JSON(http.StatusOK, jsonData)
	} else {
		data := struct {
			TotalCount int64       `json:"totalCount"`
			List       interface{} `json:"list"`
		}{
			0,
			[]interface{}{},
		}
		jsonData := lib.GetSuccess(data)
		context.JSON(http.StatusOK, jsonData)
	}

}

func GetRuntimeById(context *gin.Context) {
	var id = 0
	if val, ok := context.Get("id"); ok {
		id = val.(int)
	}
	if row, error := abtestDao.GetRuntimeById(int64(id)); error == nil {
		jsonData := lib.GetSuccess(row)
		context.JSON(http.StatusOK, jsonData)
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}

}

func AddRuntime(context *gin.Context) {
	serverName, _ := context.GetPostForm("serverName")
	strPolicyId, _ := context.GetPostForm("policyId")
	strGroupId, _ := context.GetPostForm("groupId")
	status, _ := context.GetPostForm("status")
	isDelete, _ := context.GetPostForm("isDelete")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	Status, _ := strconv.Atoi(status)
	IsDelete, _ := strconv.Atoi(isDelete)

	policyId, _ := strconv.Atoi(strPolicyId)
	groupId, _ := strconv.Atoi(strGroupId)

	item := abtestDao.Runtime{
		ServerName: serverName,
		PolicyId:   int64(policyId),
		GroupId:    int64(groupId),
		Status:     uint8(Status),
		IsDelete:   uint8(IsDelete),
		UpdateTime: time.Now().In(loc),
		CreateTime: time.Now().In(loc),
	}
	if id, err := abtestDao.AddRuntime(&item); err == nil {
		jsonData := lib.GetSuccess(id)
		context.JSON(http.StatusOK, jsonData)
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}

}

// add runtime policies online to redis
func OnlinePolicy2Redis(obj *abtestDao.Policy, responseData *[]byte) (myError lib.Error, ok bool) {
	var divData interface{}
	err := json.Unmarshal([]byte(obj.DivData), &divData)
	if err != nil {
		return lib.GetErrorByName("UNKNOWN_JSON_FORMAT"), false
	}
	postData := map[string]interface{}{
		"divtype": obj.DivModel,
		"divdata": divData,
	}
	readData, err := json.Marshal(postData)
	if err != nil {
		return lib.GetErrorByName("UNKNOWN_JSON_FORMAT"), false
	}

	*responseData, err = lib.SimpleHttpRequest("POST", *lib.TripServs["abtestPolicyOnlineServer"], string(readData), *lib.AppCfg.AppName)
	if err != nil {
		return lib.GetErrorByName("REQUEST_SERVER_ERROR"), false
	}

	return lib.GetSuccess(""), true
}

// add runtime policies online to zk
func OnlineRuntime2Zk(path string, data []byte) (myError lib.Error, ok bool) {
	conn, _, err := zk.Connect(lib.AppZkCfg.ZkServers, 3*time.Second)
	if err != nil {
		return lib.GetErrorByName("ZK_OP_ERROR"), false
	}
	ok, stat, err := conn.Exists(path)
	if !ok || err != nil {

		acl := zk.ACL{
			Perms:  0x1f, //READ | WRITE | CREATE | DELETE | ADMIN (31)
			Scheme: "world",
			ID:     "anyone",
		}
		_, err := conn.Create(path, []byte(""), 0, []zk.ACL{acl})
		if err != nil {
			return lib.GetErrorByName("ZK_OP_PATH_ERROR"), false
		}
	}
	stat, err = conn.Set(path, data, stat.Version)
	if err != nil {
		return lib.GetErrorByName("ZK_OP_PATH_ERROR"), false
	}
	logrus.Info("set "+path+" stat: ", stat)

	return lib.GetSuccess(""), true
}

// online runtime policies to redis and zk
func OnlineRuntime(Id int) (myError lib.Error, ok bool) {
	row, err := abtestDao.GetRuntimeById(int64(Id))
	if err != nil {
		return lib.GetErrorByName("DB_OP_ERROR"), false
	}
	serverName := row.ServerName

	var ids []int64
	if row.GroupId > 0 {
		groupRow, err := abtestDao.GetPolicyGroupById(row.GroupId)
		if err != nil {
			return lib.GetErrorByName("DB_OP_ERROR"), false
		}
		for _, strPolicyId := range strings.Split(groupRow.PolicyIds, ",") {
			intPolicyId, _ := strconv.Atoi(strPolicyId)
			ids = append(ids, int64(intPolicyId))
		} //end for
		policyRows, err := abtestDao.GetPolicyListByIds([]string{}, ids)
		if err != nil {
			return lib.GetErrorByName("DB_OP_ERROR"), false
		}
		for _, policyRow := range policyRows {
			//todo
			logrus.Println(policyRow)
		}
	} //end if
	if row.PolicyId > 0 {
		policyRow, err := abtestDao.GetPolicyById(row.PolicyId)
		if err != nil {
			return lib.GetErrorByName("DB_OP_ERROR"), false
		}
		path, ok := lib.AppZkCfg.NodePath[policyRow.DivModel]
		if !ok {
			return lib.GetErrorByName("ZK_NO_PATH_ERROR"), false
		}
		responseData := []byte{}
		myErr, ok := OnlinePolicy2Redis(policyRow, &responseData)
		if !ok {
			return myErr, ok
		}
		successData := map[string]interface{}{}
		buildInError := json.Unmarshal(responseData, &successData)
		if buildInError != nil {
			return lib.GetErrorByName("JSON_UNMARSHAL_ERROR"), false
		}
		code, _ := successData["code"]
		if code.(float64) != 200 {
			return lib.GetErrorByName("REQUEST_SERVER_ERROR"), false
		}
		id, _ := successData["data"]
		incrOnlinePolicyId := strconv.Itoa(int(id.(float64)))
		params := `hostname=` + serverName + `&policyid=` + incrOnlinePolicyId
		resData, err := lib.SimpleHttpRequest("GET", *lib.TripServs["abtestRuntimeOnlineServer"]+"&"+params, "", *lib.AppCfg.AppName)
		if err != nil {
			return lib.GetErrorByName("REQUEST_SERVER_ERROR"), false
		}
		buildInError = json.Unmarshal(resData, &successData)
		if buildInError != nil {
			return lib.GetErrorByName("JSON_UNMARSHAL_ERROR"), false
		}
		code, _ = successData["code"]
		if code.(float64) != 200 {
			return lib.GetErrorByName("REQUEST_SERVER_ERROR"), false
		}
		zkNodeData := []byte(`{"server_name":"` + serverName + `","policy_id":` + incrOnlinePolicyId + `}`)
		if err, ok := OnlineRuntime2Zk(path, zkNodeData); !ok {
			return err, ok
		}
	}

	return lib.GetSuccess(""), true
}

func OffRuntime(id int) (myError lib.Error, ok bool) {

	return lib.GetSuccess(""), true
}

func UpdateRuntime(context *gin.Context) {
	id, _ := context.GetPostForm("id")
	serverName, _ := context.GetPostForm("serverName")
	strPolicyId, _ := context.GetPostForm("policyId")
	strGroupId, _ := context.GetPostForm("groupId")
	status, _ := context.GetPostForm("status")
	isDelete, _ := context.GetPostForm("isDelete")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	Id, _ := strconv.Atoi(id)
	policyId, _ := strconv.Atoi(strPolicyId)
	groupId, _ := strconv.Atoi(strGroupId)
	Status, _ := strconv.Atoi(status)
	IsDelete, _ := strconv.Atoi(isDelete)

	if policyId > 0 && groupId > 0 {
		context.JSON(http.StatusOK, lib.GetErrorByName("REQUEST_PARAMS_VALIDATE_ERROR"))
		context.Abort()
		return
	}

	if Status == 3 {
		if myErr, ok := OffRuntime(Id); !ok {
			context.JSON(http.StatusOK, myErr)
			context.Abort()
			return
		}
	}
	if Status == 4 {
		if myErr, ok := OnlineRuntime(Id); !ok {
			context.JSON(http.StatusOK, myErr)
			context.Abort()
			return
		}
	}

	updatedRuntime := abtestDao.Runtime{
		Id:         int64(Id),
		ServerName: serverName,
		PolicyId:   int64(policyId),
		GroupId:    int64(groupId),
		Status:     uint8(Status),
		IsDelete:   uint8(IsDelete),
		UpdateTime: time.Now().In(loc),
	}

	updateFields := []string{"UpdateTime"}
	if len(serverName) > 0 {
		updateFields = append(updateFields, "ServerName")
	}
	if groupId > 0 && policyId == 0 {
		updateFields = append(updateFields, "GroupId")
	}
	if policyId > 0 && groupId == 0 {
		updateFields = append(updateFields, "PolicyId")
	}
	if Status >= 0 {
		updateFields = append(updateFields, "Status")
	}
	if IsDelete >= 0 {
		updateFields = append(updateFields, "IsDelete")
	}

	if err := abtestDao.UpdateRuntimeById(&updatedRuntime, updateFields...); err == nil {
		jsonData := lib.GetSuccess(id)
		context.JSON(http.StatusOK, jsonData)
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}

}

func DeleteRuntime(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err == nil && id > 0 {
		if err := abtestDao.DeleteRuntime(int64(id)); err == nil {
			jsonData := lib.GetSuccess(id)
			context.JSON(http.StatusOK, jsonData)
		} else {
			context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
		}
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("REQUEST_PARAMS_ID_ERROR"))
	}
}
