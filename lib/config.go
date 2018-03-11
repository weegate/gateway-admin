//@author wuyong
//@date   2018/1/15
//@desc easy config
//@todo add channel/lock to r/w shared config for goroutine

package lib

import (
	"flag"
	"strings"
	"sync"
)

type AppConfig struct {
	IsDev           *bool
	ServerPort      *string
	AppName         *string
	AppTpl          *AppTplConfig
	AppDb           *AppDbConfig
	AppZk           *AppZkConfig
	AppRedis        *AppRedisConfig
	TripartiteServs *TripartiteServers
	LogDir          *string
}

//todo split host uri params rpc
type TripartiteServers map[string]*string

type AppTplConfig struct {
	FeDir      *string
	LeftDelim  *string
	RightDelim *string
}

//@todo multi databases and multi table
type AppDbConfig struct {
	DbDriver   *string
	DbHost     *string
	DbPort     *string
	DbUser     *string
	DbPassword *string
	DbZone     *string
	DbName     *string
}

type AppZkConfig struct {
	ZkServers []string
	NodePath  map[string]string
	NodeData  map[string]interface{}
	RWLock    sync.RWMutex
}

type AppRedisConfig struct {
	RedisAddr     *string
	RedisPassword *string
	SelectDb      int
	KeyTpls       map[string]string
}

type AppSsoAuthConfig struct {
	AppName *string
	Version *string
	Host    *string
}

var (
	AppCfg        AppConfig
	AppTplCfg     AppTplConfig
	AppDbCfg      AppDbConfig
	AppZkCfg      AppZkConfig
	AppRedisCfg   AppRedisConfig
	AppSsoAuthCfg AppSsoAuthConfig
	TripServs     TripartiteServers
)

func init() {
	//@todo: use yaml config follow flags

	AppCfg.ServerPort = flag.String("port", ":8989", "server port")
	AppCfg.AppName = flag.String("appName", "abtest", "appName")
	AppCfg.LogDir = flag.String("logDir", "./log", "log dir path")

	AppTplCfg.FeDir = flag.String("feDir", "./fe", "frontend dir path")
	AppTplCfg.LeftDelim = flag.String("leftDelim", "{{", "left delimiter for template")
	AppTplCfg.RightDelim = flag.String("rightDelim", "}}", "right delimiter for template")

	AppDbCfg.DbDriver = flag.String("dbDriver", "mysql", "db driver")
	AppDbCfg.DbHost = flag.String("dbHost", "127.0.0.1", "db host")
	AppDbCfg.DbPort = flag.String("dbPort", "3306", "db port")
	AppDbCfg.DbUser = flag.String("dbUser", "root", "db user")
	AppDbCfg.DbPassword = flag.String("dbPassword", "root", "db password")
	AppDbCfg.DbZone = flag.String("dbZone", "Asia/Shanghai", "db zone")
	AppDbCfg.DbName = flag.String("dbName", "abtest", "db name")
	AppCfg.IsDev = flag.Bool("isDev", true, "is development env see debug log (access log,app log, db query log")

	AppRedisCfg.RedisAddr = flag.String("redisAddr", "127.0.0.1:6379", "redis address")
	AppRedisCfg.RedisPassword = flag.String("redisPassword", "foobared", "redis password")
	AppRedisCfg.SelectDb = *flag.Int("selectDb", 0, "redis select db")

	AppSsoAuthCfg.AppName = flag.String("ssoAuthAppName", *AppCfg.AppName+"-admin", "SSO auth AppName")
	AppSsoAuthCfg.Host = flag.String("ssoAuthHost", "http://sso.so", "SSO auth host")
	AppSsoAuthCfg.Version = flag.String("ssoAuthVersion", "1.0.5", "SSO auth host")

	TripServs = make(TripartiteServers)
	TripServs["abtestPolicyOnlineServer"] = flag.String("abtestPolicyOnlineServer", "http://127.0.0.1:8080/ab_admin?action=policy_set", "for abtest app policy online api address")
	TripServs["abtestPolicyGroupOnlineServer"] = flag.String("abtestPolicyGroupOnlineServer", "http://127.0.0.1:8080/ab_admin?action=policygroup_set", "for abtest app policy group online api address")
	TripServs["abtestRuntimeOnlineServer"] = flag.String("abtestRuntimeOnlineServer", "http://127.0.0.1:8080/ab_admin?action=runtime_set", "for abtest app runtime online api address")

	zkServers := flag.String("zkServers", "127.0.0.1:2181,127.0.0.1:2182,127.0.0.1:2183", "zk servers")

	divModuleNames := flag.String("divModuleName", "request_body_countryCode", "diversion module name,split by ','")
	zkNodePathNames := flag.String("zkNodePathName", "/gateway/global/runtime_policy/country_code", "zk diversion module node path name,split by ','")

	flag.Parse()

	AppZkCfg.ZkServers = strings.Split(*zkServers, ",")
	arrDivModuleNames := strings.Split(*divModuleNames, ",")
	arrZkNodePathNames := strings.Split(*zkNodePathNames, ",")
	InitNodePath(arrDivModuleNames, arrZkNodePathNames)

	InitRedisKeyConfig()

	AppCfg.AppTpl = &AppTplCfg
	AppCfg.AppDb = &AppDbCfg
	AppCfg.AppZk = &AppZkCfg
	AppCfg.AppRedis = &AppRedisCfg
	AppCfg.TripartiteServs = &TripServs
}

// init runtime policy div_model -> node path
func InitNodePath(pArrModuleNames []string, pArrNodePath []string) {
	AppZkCfg.NodePath = make(map[string]string)
	for k, v := range pArrModuleNames {
		AppZkCfg.NodePath[v] = pArrNodePath[k]
	}
}

//
func InitRedisKeyConfig() {
	switch *AppCfg.AppName {
	case "abtest":
		AppRedisCfg.KeyTpls = make(map[string]string)
		AppRedisCfg.KeyTpls["policyDivData"] = "ab:policies:{{{ID}}}:divdata"
		AppRedisCfg.KeyTpls["policyDivModel"] = "ab:policies:{{{ID}}}:divtype"
		AppRedisCfg.KeyTpls["policyIdCount"] = "ab:policies:idCount"

		AppRedisCfg.KeyTpls["policyGroup"] = "ab:policygroups:{{{GID}}}"
		AppRedisCfg.KeyTpls["policyGroupIdCount"] = "ab:policygroups:idCount"

		AppRedisCfg.KeyTpls["runtimeDivModuleName"] = "ab:runtimeInfo:{{{SERVER_NAME}}}:{{{PRIORITY}}}:divModulename"
		AppRedisCfg.KeyTpls["runtimeDivSteps"] = "ab:runtimeInfo:{{{SERVER_NAME}}}:divsteps"
		AppRedisCfg.KeyTpls["runtimeDivData"] = "ab:runtimeInfo:{{{SERVER_NAME}}}:{{{PRIORITY}}}:divDataKey"
		AppRedisCfg.KeyTpls["runtimeUserDefinedModuleName"] = "ab:runtimeInfo:{{{SERVER_NAME}}}:{{{PRIORITY}}}:userInfoModulename"
		AppRedisCfg.KeyTpls["divModuleNamePrefix"] = "abtesting.diversion"
		AppRedisCfg.KeyTpls["userDefinedModuleNamePrefix"] = "abtesting.userinfo"

		AppRedisCfg.KeyTpls["userDefinedModuleNamePrefix"] = "abtesting.userinfo"

	default:
	}

}
