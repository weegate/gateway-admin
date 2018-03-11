//@author wuyong
//@date   2018/1/15
//@desc

package abtestServerPage

import (
	"gateway-admin/model/abtest/dao"
)

func GetRutimeInfo(query map[string]string, sortBy []string, order []string, offset int64, limit int64) (res []interface{}, err error) {

	fields := []string{
		"Id",
		"ServerName",
		"PolicyId",
		"GroupId",
		"Status",
		"IsDelete",
		"UpdateTime",
		"CreateTime",
	}

	rows, err := abtestDao.GetAllRuntime(query, fields, sortBy, order, int64(offset), int64(limit))
	if err != nil {
		panic(err.Error())
	}

	var (
		groupIds        []int64
		policyIds       []int64
		mapPolicyIdName map[int64]string
		mapGroupIdName  map[int64]string
	)
	for _, row := range rows {
		mp := row.(map[string]interface{})
		if groupId, ok := mp["GroupId"]; groupId.(int64) > 0 && ok {
			groupIds = append(groupIds, groupId.(int64))
		}
		if policyId, ok := mp["PolicyId"]; policyId.(int64) > 0 && ok {
			policyIds = append(policyIds, policyId.(int64))
		}
	}

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
	}
	if len(groupIds) > 0 {
		if groupRows, err := abtestDao.GetPolicyGroupListByIds([]string{"Id", "Name"}, groupIds); err != nil {
			return nil, err
		} else {
			mapGroupIdName = make(map[int64]string, len(groupRows))
			for _, row := range groupRows {
				r := row.(map[string]interface{})
				if id, ok := r["Id"]; ok {
					if name, ok := r["Name"]; ok {
						mapGroupIdName[id.(int64)] = name.(string)
					}
				}
			} //end for
		} // end if
	}

	for _, row := range rows {
		mp := row.(map[string]interface{})
		if groupId, ok := mp["GroupId"]; groupId.(int64) > 0 && ok {
			mp["GroupName"] = mapGroupIdName[groupId.(int64)]
		}
		if policyId, ok := mp["PolicyId"]; policyId.(int64) > 0 && ok {
			mp["PolicyName"] = mapPolicyIdName[policyId.(int64)]
		}
		res = append(res, mp)
	}

	return
}
