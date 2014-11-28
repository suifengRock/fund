export GOPATH := $(shell pwd)
export PATH := ${PATH}:${GOPATH}\bin
export GOBIN := ${GOPATH}\bin
DOCKER_RUN_GO := fig run --rm goapp
FIG_VERSION := $(shell fig --version 2>/dev/null)

main:clear-pkg
	go run main.go
build:
	go install main.go
images:
	docker build -t fund .
run:
	docker run -it -v /Users/suifengluo/Documents/gospace/fund:/gopath --rm fund
docker: images run

dependence:
	go get github.com/go-sql-driver/mysql
	go get github.com/go-xorm/xorm
	go get github.com/qiniu/iconv
	go get github.com/PuerkitoBio/goquery

fig: 
ifdef FIG_VERSION
	@echo "Found fig version $(FIG_VERSION)"
else
	@echo fig Not found try to install it
	curl -L https://github.com/docker/fig/releases/download/1.0.1/fig-`uname -s`-`uname -m` > /usr/local/bin/fig; chmod +x /usr/local/bin/fig
endif


clear-pkg:
	rm -rf pkg/

shell: fig
	$(DOCKER_RUN_GO) bash

mysqlc: fig
	$(DOCKER_RUN_GO) go run mysql.go

mainc: fig
	$(DOCKER_RUN_GO) go run main.go