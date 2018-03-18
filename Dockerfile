FROM golang:1.10.0-alpine AS build

RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/github.com/orisano/preve/
WORKDIR /go/src/github.com/orisano/preve
RUN dep ensure -vendor-only

COPY . /go/src/github.com/orisano/preve
RUN go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags '-static'" -o bin/check github.com/orisano/preve/cmd/check
RUN go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags '-static'" -o bin/in github.com/orisano/preve/cmd/in

FROM scratch
COPY --from=build /go/src/github.com/orisano/preve/bin/check /opt/resource/check
COPY --from=build /go/src/github.com/orisano/preve/bin/in /opt/resource/in
