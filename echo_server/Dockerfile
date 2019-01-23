FROM golang:1.11
WORKDIR /go/src/github.com/pmoncadaisla/istio-auth-sample/echo_server/
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep ensure -v -update
RUN CGO_ENABLED=0 go build -o app

FROM alpine:latest  
RUN apk --no-cache add ca-certificates curl
WORKDIR /
COPY --from=0 /go/src/github.com/pmoncadaisla/istio-auth-sample/echo_server/app .
ENTRYPOINT ["./app"]  