package abtestDao

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PolicyGroup struct {
	Id         int64     `orm:"column(id);auto"`
	Name       string    `orm:"column(name);size(255)" description:"组名"`
	PolicyIds  string    `orm:"column(policy_ids);size(255)" description:"策略id，最多10个，用,分割"`
	UpdateTime time.Time `orm:"column(update_time);type(datetime)" description:"更新时间"`
	CreateTime time.Time `orm:"column(create_time);type(datetime)" description:"创建时间"`
	Status     uint8     `orm:"column(status)" description:"状态：0.通过，1.审核中，2.拒绝"`
	IsDelete   uint8     `orm:"column(is_delete)" description:"状态：0.有效，1.删除"`
	Ext1       string    `orm:"column(ext1);size(255)"`
	Ext2       uint      `orm:"column(ext2)"`
}

func (t *PolicyGroup) TableName() string {
	return "policy_group"
}

func init() {
	orm.RegisterModel(new(PolicyGroup))
}

// AddPolicyGroup insert a new PolicyGroup into database and returns
// last inserted Id on success.
func AddPolicyGroup(m *PolicyGroup) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPolicyGroupById retrieves PolicyGroup by Id. Returns error if
// Id doesn't exist
func GetPolicyGroupById(id int64) (v *PolicyGroup, err error) {
	o := orm.NewOrm()
	v = &PolicyGroup{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetPolicyGroupListByIds(fields []string, ids []int64) (res []interface{}, err error) {
	var rows []PolicyGroup

	if _, err = orm.NewOrm().QueryTable(new(PolicyGroup)).Filter("id__in", ids).All(&rows, fields...); err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		for _, v := range rows {
			res = append(res, v)
		}
	} else {
		// trim unused fields
		for _, v := range rows {
			m := make(map[string]interface{})
			val := reflect.ValueOf(v)
			for _, fieldName := range fields {
				m[fieldName] = val.FieldByName(fieldName).Interface()
			}
			res = append(res, m)
		}
	}

	return res, err
}

// GetAllPolicyGroup retrieves all PolicyGroup matches certain condition. Returns empty list if
// no records exist
func GetAllPolicyGroup(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PolicyGroup))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []PolicyGroup
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdatePolicyGroup updates PolicyGroup by Id and returns error if
// the record to be updated doesn't exist
func UpdatePolicyGroupById(m *PolicyGroup, cols ...string) (err error) {
	o := orm.NewOrm()
	v := PolicyGroup{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePolicyGroup deletes PolicyGroup by Id and returns error if
// the record to be deleted doesn't exist
func DeletePolicyGroup(id int64) (err error) {
	o := orm.NewOrm()
	v := PolicyGroup{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PolicyGroup{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetPolicyGroupTotalNum() (int64, error) {
	o := orm.NewOrm()
	total, err := o.QueryTable("policy_group").Count()
	if err != nil {
		return 0, errors.New("sql execute failed")
	}

	return total, nil
}
