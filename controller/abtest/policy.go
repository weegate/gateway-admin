//@author wuyong
//@date   2018/1/9
//@desc

package abtestController

import (
	"net/http"
	"strconv"

	"gateway-admin/lib"
	_ "gateway-admin/model"
	"time"

	"gateway-admin/model/abtest/dao"

	"github.com/gin-gonic/gin"
)

func GetPolicyList(context *gin.Context) {
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

	fields := []string{
		"Id",
		"Name",
		"UpdateTime",
		"CreateTime",
		"Status",
		"IsDelete",
		"DivModel",
		"DivData",
	}
	sortBy := []string{
		"create_time",
	}
	order := []string{
		"desc",
	}

	rows, err := abtestDao.GetAllPolicy(query, fields, sortBy, order, int64(offset), int64(limit))
	if err != nil {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}
	total, err := abtestDao.GetPolicyTotalNum()
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

func GetPolicyById(context *gin.Context) {
	var id = 0
	if val, ok := context.Get("id"); ok {
		id = val.(int)
	}
	if row, error := abtestDao.GetPolicyById(int64(id)); error == nil {
		jsonData := lib.GetSuccess(row)
		context.JSON(http.StatusOK, jsonData)
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}
}

func AddPolicy(context *gin.Context) {
	name, _ := context.GetPostForm("name")
	divModel, _ := context.GetPostForm("divmodel")
	divData, _ := context.GetPostForm("divdata")
	status, _ := context.GetPostForm("status")
	isDelete, _ := context.GetPostForm("is_delete")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	Status, _ := strconv.Atoi(status)
	IsDelete, _ := strconv.Atoi(isDelete)

	policy := abtestDao.Policy{
		Name:       name,
		DivModel:   divModel,
		DivData:    divData,
		Status:     uint8(Status),
		IsDelete:   uint8(IsDelete),
		UpdateTime: time.Now().In(loc),
		CreateTime: time.Now().In(loc),
	}
	if id, err := abtestDao.AddPolicy(&policy); err == nil {
		jsonData := lib.GetSuccess(id)
		context.JSON(http.StatusOK, jsonData)
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}

}

func UpdatePolicy(context *gin.Context) {
	id, _ := context.GetPostForm("id")
	name, _ := context.GetPostForm("name")
	divModel, _ := context.GetPostForm("divmodel")
	divData, _ := context.GetPostForm("divdata")
	status, _ := context.GetPostForm("status")
	isDelete, _ := context.GetPostForm("is_delete")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	Id, _ := strconv.Atoi(id)
	Status, _ := strconv.Atoi(status)
	IsDelete, _ := strconv.Atoi(isDelete)

	updatedPolicy := abtestDao.Policy{
		Id:         int64(Id),
		Name:       name,
		DivModel:   divModel,
		DivData:    divData,
		Status:     uint8(Status),
		IsDelete:   uint8(IsDelete),
		UpdateTime: time.Now().In(loc),
	}
	updateFields := []string{"Name", "DivModel", "DivData", "Status", "IsDelete", "UpdateTime"}

	if err := abtestDao.UpdatePolicyById(&updatedPolicy, updateFields...); err == nil {
		jsonData := lib.GetSuccess(id)
		context.JSON(http.StatusOK, jsonData)
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
	}

}

func DeletePolicy(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err == nil && id > 0 {
		if err := abtestDao.DeletePolicy(int64(id)); err == nil {
			jsonData := lib.GetSuccess(id)
			context.JSON(http.StatusOK, jsonData)
		} else {
			context.JSON(http.StatusOK, lib.GetErrorByName("DB_OP_ERROR"))
		}
	} else {
		context.JSON(http.StatusOK, lib.GetErrorByName("REQUEST_PARAMS_ID_ERROR"))
	}
}
