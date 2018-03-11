package abtestDao

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Runtime struct {
	Id         int64     `orm:"column(id);auto"`
	ServerName string    `orm:"column(server_name);size(255)" description:"运行时策略对应的服务名"`
	PolicyId   int64     `orm:"column(policy_id)"`
	GroupId    int64     `orm:"column(group_id)"`
	Status     uint8     `orm:"column(status)" description:"状态：0.通过，1.审核中，2.拒绝，3.下线，4.上线"`
	IsDelete   uint8     `orm:"column(is_delete)" description:"状态：0.可用，1.删除"`
	UpdateTime time.Time `orm:"column(update_time);type(datetime)" description:"更新时间"`
	CreateTime time.Time `orm:"column(create_time);type(datetime)" description:"创建时间"`
	Ext1       string    `orm:"column(ext1);size(255)"`
	Ext2       uint      `orm:"column(ext2)"`
}

func (t *Runtime) TableName() string {
	return "runtime"
}

func init() {
	orm.RegisterModel(new(Runtime))
}

// AddRuntime insert a new Runtime into database and returns
// last inserted Id on success.
func AddRuntime(m *Runtime) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRuntimeById retrieves Runtime by Id. Returns error if
// Id doesn't exist
func GetRuntimeById(id int64) (v *Runtime, err error) {
	o := orm.NewOrm()
	v = &Runtime{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRuntime retrieves all Runtime matches certain condition. Returns empty list if
// no records exist
func GetAllRuntime(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Runtime))
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

	var l []Runtime
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

// UpdateRuntime updates Runtime by Id and returns error if
// the record to be updated doesn't exist
func UpdateRuntimeById(m *Runtime, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Runtime{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRuntime deletes Runtime by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRuntime(id int64) (err error) {
	o := orm.NewOrm()
	v := Runtime{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Runtime{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetRuntimeTotalNum() (int64, error) {
	o := orm.NewOrm()
	total, err := o.QueryTable("runtime").Count()
	if err != nil {
		return 0, errors.New("sql execute failed")
	}

	return total, nil
}
