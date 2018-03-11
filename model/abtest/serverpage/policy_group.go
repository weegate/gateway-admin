//@author wuyong
//@date   2018/1/15
//@desc

package abtestServerPage

import (
	"gateway-admin/model/abtest/dao"
	"strconv"
	"strings"
)

func GetPolicyGroupInfo(query map[string]string, sortBy []string, order []string, offset int64, limit int64) (res []interface{}, err error) {

	fields := []string{
		"Id",
		"Name",
		"PolicyIds",
		"Status",
		"IsDelete",
		"UpdateTime",
		"CreateTime",
	}

	rows, err := abtestDao.GetAllPolicyGroup(query, fields, sortBy, order, int64(offset), int64(limit))
	if err != nil {
		panic(err.Error())
	}

	var (
		policyIds       []int64
		mapPolicyIdName map[int64]string
	)
	for _, row := range rows {
		mp := row.(map[string]interface{})
		if strPolicyIds, ok := mp["PolicyIds"]; len(strPolicyIds.(string)) > 0 && ok {
			for _, id := range strings.Split(strPolicyIds.(string), ",") {
				if policyId, err := strconv.Atoi(id); err == nil {
					policyIds = append(policyIds, int64(policyId))
				}
			} //end for
		} //edn if
	} //end for

	if len(policyIds) > 0 {
		if policyRows, err := abtestDao.GetPolicyListByIds([]string{"Id", "Name"}, policyIds); err != nil {
			return nil, err
		} else {
			mapPolicyIdName = make(map[int64]string, len(policyRows))
			for _, row := range policyRows {
				r := row.(map[string]interface{})
				if id, ok := r["Id"]; ok {
					if name, ok := r["Name"]; ok {
						mapPolicyIdName[id.(int64)] = name.(string)
					}
				}
			} //end for
		} //end if
	} //end if

	for _, row := range rows {
		mp := row.(map[string]interface{})
		var strPolicies []string
		if strPolicyIds, ok := mp["PolicyIds"]; len(strPolicyIds.(string)) > 0 && ok {
			for _, id := range strings.Split(strPolicyIds.(string), ",") {
				if policyId, err := strconv.Atoi(id); err == nil && policyId > 0 {
					strPolicies = append(strPolicies, mapPolicyIdName[int64(policyId)])
				}
			} //end for
			mp["Policies"] = strings.Join(strPolicies, ",")
		} //end if
		res = append(res, mp)
	} // end for

	return
}
