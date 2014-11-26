export GOPATH :=F:\working\gopro\fund
export PATH := ${PATH}:${GOPATH}\bin
export GOBIN := ${GOPATH}\bin
main:
	go run main.go
build:
	go install main.go
images:
	docker build -t fund .
run:
	docker run -it -v F:\working\gopro\fund:/gopath --rm fund
docker: images run
