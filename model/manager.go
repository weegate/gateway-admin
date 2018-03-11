//@author wuyong
//@date   2018/1/10
//@desc todo: manage the db(master(w)/slave(r)) and init by app model  more....do multi databases and multi tables by using zk

package model

import (
	"net/url"

	"fmt"

	"gateway-admin/lib"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
)

func init() {
	if *lib.AppCfg.IsDev {
		orm.Debug = true
	}

	dataSource := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&loc=%v", *lib.AppDbCfg.DbUser, *lib.AppDbCfg.DbPassword, *lib.AppDbCfg.DbHost, *lib.AppDbCfg.DbPort, *lib.AppDbCfg.DbName, url.QueryEscape(*lib.AppDbCfg.DbZone))

	//todo change it use zk to manage multi database and multi table -> make a mw
	err := orm.RegisterDataBase("default", *lib.AppDbCfg.DbDriver, dataSource, 30)
	if err != nil {
		panic(err.Error())
	}

	// create table
	//orm.RunSyncdb("default", false, true)
}

func manageDb(app string, module string) {

}
