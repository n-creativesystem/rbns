NAME	:= rbns
SRCS	:= $(shell find . -type d -name archive -prune -o -type f -name '*.go')
LDFLAGS	:= -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

build/frontend:
	@yarn build

build/static: $(SRCS)
	CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

build: $(SRCS)
	go build -o bin/$(NAME)

.PHONY: deps
deps:
	go get -v

.PHONY: cross-build
cross-build: deps
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done

protoc:
	@protoc -I ./proto/docs --go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go_out=./proto --go_opt=paths=source_relative ./proto/docs/*.proto
	@protoc -I ./proto/docs --doc_out=./proto/docs/resource/custom_markdown.tpl,index.md:./proto/docs ./proto/docs/*.proto

.PHONY: ssl
ssl:
	@openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout ssl/server.key -out ssl/server.crt -subj "/C=JP/ST=Osaka/L=Osaka/O=NCreativeSystem, Inc./CN=localhost"

all-in-one: docker/all-in-one.dockerfile docker/*.docker
	cpp -P -o all-in-one.dockerfile docker/all-in-one.dockerfile
backend: docker/api-only.dockerfile docker/*.docker
	cpp -P -o backend.dockerfile docker/api-only.dockerfile

build/all-in-one: all-in-one
	docker build -t ${IMAGE_NAME} -f all-in-one.dockerfile .
build/backend: backend
	docker build -t ${IMAGE_NAME} -f backend.dockerfile .