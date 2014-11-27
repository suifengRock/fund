export GOPATH := F:\working\fund
export PATH := ${PATH}:${GOPATH}\bin
export GOBIN := ${GOPATH}\bin
main:
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
	go get github.com/PuerkitoBio/goquery


ssh:
	ssh -i cloud.key ec2_user@<instance_ip>


utf:
	go get -u github.com/djimenez/iconv-go

