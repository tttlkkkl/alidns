FROM golang:1.12.4-alpine AS build_deps

RUN apk add --no-cache git

WORKDIR /workspace
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o webhook -ldflags '-s -w -extldflags "-static"' .

FROM alpine:3.10
ENV ALICLOUD_ACCESS_KEY=
ENV ALICLOUD_SECRET_KEY=
ENV REGIONID=
RUN apk add --no-cache ca-certificates

COPY --from=build /workspace/webhook /usr/local/bin/webhook

ENTRYPOINT [ "webhook" ]