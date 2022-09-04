FROM golang:1.19

WORKDIR /go/src/github.com/ryoshindo/envoy-control-plane

COPY go.mod go.sum ./

RUN go install github.com/cosmtrek/air
