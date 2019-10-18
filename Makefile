APP?=app
PORT?=1323
RELEASE?=0.0.1
#COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROJECT?=mq/academy
GOOS?=linux
GOARCH?=amd64
REGISTRY_ADDR?=127.0.0.1:5000

clean:
	rm -f ${APP}

base:
	docker build -t $(APP)-base-build -f Dockerfile-build-base .

build: clean
#	go get -v ./...
	CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build \
	-mod vendor \
	-ldflags "-s -w -X main.Release=${RELEASE} \
	-X main.BuildTime=${BUILD_TIME}" \
	-o ${APP}

run: image
#	SERVER_PORT=${PORT} ./${APP}
	docker stop $(APP) || true
	docker run --name ${APP} -p ${PORT}:${PORT} --rm \
		-e "SERVER_PORT=${PORT}" \
		$(REGISTRY_ADDR)/$(APP):$(RELEASE)

test:
	go test -v -race -cover ./...

image:
	docker build -t $(REGISTRY_ADDR)/$(APP):$(RELEASE) .

push:
	docker push $(REGISTRY_ADDR)/$(APP):$(RELEASE)