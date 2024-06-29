.PHONY: build

# set local variables for function name to be substituted in the template
FN ?= RobloxTokenAuthorizer


build-sam:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main . && \
	cp ./main ./.aws-sam/build/${FN}/main && \
	cp -R ./templates/ .aws-sam/build/${FN}/
	sam build && sam deploy


build-${FN}:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main && \
	cp ./main ./.aws-sam/build/${FN}/main

build-dlv:
	go get github.com/go-delve/delve/cmd/dlv
	GOOS=linux GOARCH=amd64 go build -o ./debugger/dlv github.com/go-delve/delve/cmd/dlv

launch-invoke:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main && \
	cp ./main ./.aws-sam/build/${FN}/main && \
	cp ./debugger/dlv ./.aws-sam/build/${FN}/dlv && \
	cp -R templates/ ./.aws-sam/build/${FN}/ && \
	sam local invoke ${FN} --event ./testinputs/authevent.json

launch-debug:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main && \
	cp ./main ./.aws-sam/build/${FN}/main && \
	cp ./debugger/dlv ./.aws-sam/build/${FN}/dlv && \
	sam build && \
	sam local invoke ${FN} --event ./authevent.json --debug-port 12312 --debug-args "-delveAPI=2" --debugger-path ./debugger

launch-api:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main && \
	cp ./main ./.aws-sam/build/${FN}/main && \
	cp ./debugger/dlv ./.aws-sam/build/${FN}/dlv && \
	cp -R templates/ ./.aws-sam/build/${FN}/ && \
	sam local start-api
	
initdb:
	go run ddl/init.go
