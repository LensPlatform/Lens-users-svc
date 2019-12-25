# Makefile for releasing lens platform service
#
# The release version is controlled from pkg/version
TAG?=latest
NAME:=Lens-users-svc
DOCKER_REPOSITORY:=lensPlatform
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
GIT_COMMIT:=$(shell git describe --dirty --always)
VERSION:=$(shell grep 'VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')

generate-proto:
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/infobloxopen/protoc-gen-gorm
	go get -u github.com/infobloxopen/atlas-app-toolkit
	protoc -I/usr/local/include -I. \
			-I$(GOPATH)/src \
			--go_out=plugins=grpc:./pkg/model/ \
			--grpc-gateway_out=logtostderr=true:. \
			--gorm_out="engine=postgres:." ./proto/common/*.proto
	protoc -I. \
			-I=vendor/github.com/lyft/protoc-gen-validate \
			-I=vendor/github.com/infobloxopen/protoc-gen-gorm \
			-I=vendor/github.com/infobloxopen/atlas-app-toolkit \
			--go_out=plugins=grpc:./pkg/model/ \
			--grpc-gateway_out=logtostderr=true:.
			--gorm_out="engine=postgres:." ./proto/investor-user/*.proto
	protoc -I/usr/local/include \
			-I. \
			-I$(GOPATH)/src \
			-I=vendor/github.com/grpc-ecosystem/grpc-gateway/ \
			-I=vendor/github.com/infobloxopen/atlas-app-toolkit/rpc/resource/ \
			-I=vendor/github.com/infobloxopen/protoc-gen-gorm/options/ \
			--go_out=plugins=grpc:./pkg/model/ --grpc-gateway_out=logtostderr=true:. \
			--gorm_out="engine=postgres:." ./proto/group/*.proto
	protoc --go_out=plugins=grpc:./pkg/model/ --grpc-gateway_out=logtostderr=true:. --gorm_out="engine=postgres:." ./proto/startup-user/*.proto
	protoc --go_out=plugins=grpc:./pkg/model/ --grpc-gateway_out=logtostderr=true:. --gorm_out="engine=postgres:." ./proto/team/*.proto
	protoc --go_out=plugins=grpc:./pkg/model/ --grpc-gateway_out=logtostderr=true:. --gorm_out="engine=postgres:." ./proto/user/*.proto
	protoc --go_out=plugins=grpc:./pkg/model/ --grpc-gateway_out=logtostderr=true:. --gorm_out="engine=postgres:." ./proto/service.proto

clean:
	GO111MODULE= ./scripts/cleanup.sh && cd ..

format:
	gofmt -s -w .

start-services:
	docker-compose up

run:
	GO111MODULE=on go run -ldflags "-s -w -X github.com/$(DOCKER_REPOSITORY)/$(NAME)/pkg/version.REVISION=$(GIT_COMMIT)" cmd/podinfo/* \
	--level=debug --grpc-port=9999 --backend-url=https://httpbin.org/status/401 --backend-url=https://httpbin.org/status/500 \
	--ui-logo=https://raw.githubusercontent.com/stefanprodan/podinfo/gh-pages/cuddle_clap.gif --ui-color=#34577c

test:
	GO111MODULE=on go test -v -race ./...

build:
	GO111MODULE=on GIT_COMMIT=$$(git rev-list -1 HEAD) && GO111MODULE=on CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/$(DOCKER_REPOSITORY)/$(NAME)/pkg/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/podinfo ./cmd/podinfo/*
	GO111MODULE=on GIT_COMMIT=$$(git rev-list -1 HEAD) && GO111MODULE=on CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/$(DOCKER_REPOSITORY)/$(NAME)/pkg/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/podcli ./cmd/podcli/*

build-charts:
	helm lint charts/*
	helm package charts/*

build-container:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

test-container:
	@docker rm -f $(NAME) || true
	@docker run -dp 9898:9898 --name=$(NAME) $(DOCKER_IMAGE_NAME):$(VERSION)
	@docker ps
	@TOKEN=$$(curl -sd 'test' localhost:9898/token | jq -r .token) && \
	curl -sH "Authorization: Bearer $${TOKEN}" localhost:9898/token/validate | grep test

push-container:
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_IMAGE_NAME):latest
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) quay.io/$(DOCKER_IMAGE_NAME):$(VERSION)
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) quay.io/$(DOCKER_IMAGE_NAME):latest
	docker push quay.io/$(DOCKER_IMAGE_NAME):$(VERSION)
	docker push quay.io/$(DOCKER_IMAGE_NAME):latest

version-set:
	@next="$(TAG)" && \
	current="$(VERSION)" && \
	sed -i '' "s/$$current/$$next/g" pkg/version/version.go && \
	sed -i '' "s/tag: $$current/tag: $$next/g" charts/podinfo/values.yaml && \
	sed -i '' "s/appVersion: $$current/appVersion: $$next/g" charts/podinfo/Chart.yaml && \
	sed -i '' "s/version: $$current/version: $$next/g" charts/podinfo/Chart.yaml && \
	sed -i '' "s/podinfo:$$current/podinfo:$$next/g" kustomize/deployment.yaml && \
	echo "Version $$next set in code, deployment, chart and kustomize"

release:
	git tag $(VERSION)
	git push origin $(VERSION)

swagger:
	GO111MODULE=on go get github.com/swaggo/swag/cmd/swag
	cd pkg/api && $$(go env GOPATH)/bin/swag init -g server.go