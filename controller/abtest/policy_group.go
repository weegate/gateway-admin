//@author wuyong
//@date   2018/1/11
//@desc

package abtestController

import (
	"net/http"
	"strconv"
	"time"

	"gateway-admin/lib"
	_ "gateway-admin/model"
	"gateway-admin/model/abtest/dao"

	"gateway-admin/model/abtest/serverpage"

	"github.com/gin-gonic/gin"
)

func GetPolicyGroupList(context *gin.Context) {
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

	rows, err := abtestServerPage.GetPolicyGroupInfo(query, sortBy, order, int64(offset), int64(limit))
	if err != nil {
		panic(err.Error())
	}
	total, err := abtestDao.GetPolicyGroupTotalNum()
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

func GetPolicyGroupById(context *gin.Context) {
	var id = 0
	if val, ok := context.Get("id"); ok {
		id = val.(int)
	}
	if row, error := abtestDao.GetPolicyGroupById(int64(id)); error == nil {
		jsonData := lib.GetSuccess(row)
		context.JSON(http.StatusOK, jsonData)
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}

}
func AddPolicyGroup(context *gin.Context) {
	name, _ := context.GetPostForm("name")
	policyIds, _ := context.GetPostForm("policyIds")
	status, _ := context.GetPostForm("status")
	isDelete, _ := context.GetPostForm("isDelete")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	Status, _ := strconv.Atoi(status)
	IsDelete, _ := strconv.Atoi(isDelete)

	item := abtestDao.PolicyGroup{
		Name:       name,
		PolicyIds:  policyIds,
		Status:     uint8(Status),
		IsDelete:   uint8(IsDelete),
		UpdateTime: time.Now().In(loc),
		CreateTime: time.Now().In(loc),
	}
	if id, err := abtestDao.AddPolicyGroup(&item); err == nil {
		jsonData := lib.GetSuccess(id)
		context.JSON(http.StatusOK, jsonData)
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}

}

func UpdatePolicyGroup(context *gin.Context) {
	id, _ := context.GetPostForm("id")
	name, _ := context.GetPostForm("name")
	policyIds, _ := context.GetPostForm("policyIds")
	status, _ := context.GetPostForm("status")
	isDelete, _ := context.GetPostForm("isDelete")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	Id, _ := strconv.Atoi(id)
	Status, _ := strconv.Atoi(status)
	IsDelete, _ := strconv.Atoi(isDelete)

	updatedPolicyGroup := abtestDao.PolicyGroup{
		Id:         int64(Id),
		Name:       name,
		PolicyIds:  policyIds,
		Status:     uint8(Status),
		IsDelete:   uint8(IsDelete),
		UpdateTime: time.Now().In(loc),
	}
	updateFields := []string{"Name", "PolicyIds", "Status", "IsDelete", "UpdateTime"}

	if err := abtestDao.UpdatePolicyGroupById(&updatedPolicyGroup, updateFields...); err == nil {
		jsonData := lib.GetSuccess(id)
		context.JSON(http.StatusOK, jsonData)
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}

}

func DeletePolicyGroup(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err == nil && id > 0 {
		if err := abtestDao.DeletePolicyGroup(int64(id)); err == nil {
			jsonData := lib.GetSuccess(id)
			context.JSON(http.StatusOK, jsonData)
		} else {
			context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
		}
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("REQUEST_PARAMS_ID_ERROR"))
	}
}
