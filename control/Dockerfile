FROM golang:1.20

WORKDIR /go/src/github.com/ryoshindo/envoy-control-plane/control

COPY go.mod go.sum ./

RUN go install github.com/cosmtrek/air
