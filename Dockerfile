FROM golang:1.13-alpine3.10 AS build

WORKDIR /go/src/github.com/orisano/preve
COPY . .
RUN go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags '-static'" -o bin/check github.com/orisano/preve/cmd/check
RUN go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags '-static'" -o bin/in github.com/orisano/preve/cmd/in

FROM scratch
COPY --from=build /go/src/github.com/orisano/preve/bin/* /opt/resource/
