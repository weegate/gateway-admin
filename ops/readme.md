#### 部署
以下是qa环境,每个idc机房单独采用一台部署整套服务；product环境可以将以下4个服务单独部署，各个服务在此基础上进行机器扩展 (k8s->docker)

> lbs (ngx) (maybe upstream backend service is fixed)
- Cloud03 
	- 10.3.26.49
- Cloud04 
	- 10.4.22.148
- Cloud05
	- 10.5.27.237

> gateway(ngx+lua/go)  + update local cache(lmdb) script worker(python) + online runtime policy service(ngx+lua/go->redis) *consumer*
- Cloud03 
	- 10.3.26.49
- Cloud04 
	- 10.4.22.148
- Cloud05
	- 10.5.27.237

> app api service(nodejs/go) *provider*
- Cloud03 
	- 10.3.26.49
- Cloud04 
	- 10.4.22.148
- Cloud05
	- 10.5.27.237

> abtest gateway admin (go)  部署在cloud05上(相同的时区);将运行时策略分发到不同的idc配置服务中 *configer*
- Cloud05
	- 10.5.27.238

