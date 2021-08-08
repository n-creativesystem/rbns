FROM golang:1.16-alpine as build
ENV TZ=Asia/Tokyo
WORKDIR /src/
COPY go.mod go.sum Makefile ./
RUN apk update && apk add make git && go mod download
COPY auth auth
COPY config config
COPY di di
COPY domain domain
COPY handler handler
COPY infra infra
COPY logger logger
COPY proto proto
COPY protoconv protoconv
COPY service service
COPY utilsconv utilsconv
COPY *.go ./
RUN make build
FROM node:16-alpine as frontend
WORKDIR /src/
COPY frontend frontend
COPY public public
COPY package.json vue.config.js yarn.lock ./
RUN yarn global add @vue/cli && yarn install && yarn build
FROM alpine:3
RUN addgroup -g 70 -S api-rback && adduser -u 70 -S -D -G api-rback -H -h /var/lib/api-rback -s /bin/sh api-rback && mkdir -p /var/lib/api-rback && chown -R api-rback:api-rback /var/lib/api-rback
WORKDIR /var/lib/api-rback
COPY --from=frontend /src/static static/
COPY --from=build --chown=api-rback:api-rback /src/bin/rbns .
RUN chmod +x rbns && mv rbns /usr/local/bin/
USER api-rback
CMD [ "rbns" ]
