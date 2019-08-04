FROM golang:1.12-alpine3.10 AS build

WORKDIR /go/src/github.com/orisano/preve
RUN wget -O /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep
RUN apk add --no-cache git
COPY Gopkg.lock Gopkg.toml .
RUN dep ensure -vendor-only && rm -rf /go/pkg/dep/sources

COPY . .
RUN go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags '-static'" -o bin/check github.com/orisano/preve/cmd/check
RUN go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags '-static'" -o bin/in github.com/orisano/preve/cmd/in

FROM scratch
COPY --from=build /go/src/github.com/orisano/preve/bin/* /opt/resource/
