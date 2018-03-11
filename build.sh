#!/usr/bin/env bash

build(){
	appname="abtest"
	[ "x$1" != "x" ] && appname=$1
	set -v on
	mkdir -p bin
	go build -o bin/${appname}-gateway-admin.dev
	GOOS=linux go build -o bin/${appname}-gateway-admin
	set -v off
}

main(){
	if [[ $1 == "init" ]];then
		return
	else
		build $1
	fi
}

main $1
