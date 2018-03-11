#!/usr/bin/env bash
fe_path=/data/fe
sudo mkdir -p $fe_path
sudo chown -R www:www $fe_path

log_path=/data/logs/gateway/abtest-admin
mkdir -p $log_path

git clone https://github.com/weegate/gateway-admin-fe.git ${fe_path}/gateway-admin

git clone https://github.com/weegate/gateway-admin.git
cd gateway-admin && sh build.sh



# notice: before run this, must create database by mysql db (apollo)
# for local dev
if [ "x$NODE_ENV" == "xdev" ];then
	bin/abtest-gateway-admin.dev \
		-appName="abtest" \
		-port=":8989" \
		-abtestPolicyGroupOnlineServer="http://127.0.0.1:8080/ab_admin?action=policygroup_set" \
		-abtestPolicyOnlineServer="http://127.0.0.1:8080/ab_admin?action=policy_set" \
		-abtestRuntimeOnlineServer="http://127.0.0.1:8080/ab_admin?action=runtime_set" \
		-dbDriver="mysql" \
		-dbHost="127.0.0.1" \
		-dbUser="root" \
		-dbPassword="root" \
		-dbPort="3306" \
		-dbName="abtest" \
		-dbZone="Asia/Shanghai" \
		-feDir="${fe_path/gateway-admin}" \
		-isDev=true \
		-leftDelim="{{" \
		-rightDelim="}}" \
		-logDir="${log_path}" \
		-redisAddr="127.0.0.1:6379" \
		-redisPassword="foobared" \
		-selectDb=0 \
		-divModuleName="request_body_countryCode" \
		-zkNodePathName="/gateway/global/runtime_policy/country_code" \
		--zkServers="127.0.0.1:2181,127.0.0.1:2182,127.0.0.1:2183" \
		-ssoAuthHost="http://sso.in.test.so" \
		-ssoAuthVersion="1.0.5" \
		-ssoAuthAppName="test-global-abtest-admin"
		
fi



