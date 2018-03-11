### intro
	对应灰度分流策略网关，后台管理策略平台，编辑发布策略，审核通过上线/下线。

### 依赖
>开发库：  
1. beego orm (mysql)  
	 (tips: 安装bee工具来直接初始生成数据库表对应操作的model层dao文件)
	 (eg:`bee generate appcode -tables="policy,policy_group,runtime" -driver=mysql -conn="root:root@tcp(127.0.0.1:3306)/abtest" -level=1`)
2. gin Enginer (context,render,httprouter; mw -> https://github.com/gin-gonic/contrib) 
	
>服务：  
1. 网关后台api服务: 将运行时策略写入redis中,作为远程缓存
2. 配置服务: 将运行时策略元数据写入已开发上线策略模块的节点中 -> 网关本地异步脚本实时watch更新本地缓存
	

### 目录结构
````
├── README.md
├── auth
│   ├── sso.go
│   └── sso_test.go
├── bin
├── controller
│   ├── abtest
│   │   ├── policy.go
│   │   ├── policy_group.go
│   │   └── runtime.go
│   └── dispatcher.go
├── fe
│   ├── mock
│   ├── static
│   │   ├── bootstrap
│   │   ├── css
│   │   ├── img
│   │   └── js
│   └── template
│       ├── abtest
│       └── tpl.index.html
├── glide.lock
├── glide.yaml
├── lib
│   ├── config.go
│   ├── error.go
│   ├── session.go
│   └── util.go
├── log
│   ├── access.log
│   └── app.log
├── main.go
├── model
│   ├── abtest
│   │   ├── dao
│   │   └── serverpage
│   └── manager.go
├── mw
│   └── auth.go
├── vendor
│   └── 
└── view
    ├── abtest
    │   ├── policy.go
    │   ├── policy_group.go
    │   └── runtime.go
    └── render.go
````

### howto
> - 定义运行配置：运行脚本`sh build.sh abtest` 通过 -h 查看对应的参数配置运行`bin/abtest-gateway-admin.dev -h`
> - 部署在cloud05相关的机器上，相同的zone

### todo
- [ ] 0.RBAC (审核机制)  
- [ ] 1.采用[json-editor](http://jeremydorn.com/json-editor/)对策略schema进行编辑
- [ ] 2.选择的运行时策略定时上线生效
- [ ] 3.根据时间将日志切割logrotate
- [ ] 4.OPDYPM(open platform diy your policy module)添加用户已开发好的分流模块管理UI(开发完某个新场景的分流模块，测试通过确认上线后，扫描对应的分流模块目录和用户自定义的解析模块目录中的文件名写入表中)
- [ ] _异步worker脚本实时watch监控线上运行策略情况(不属于这个后台，辅助功能脚本worker)_
