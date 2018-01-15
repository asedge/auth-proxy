FROM golang:latest
COPY . /go/src/github.com/asedge/auth-proxy
WORKDIR /go/src/github.com/asedge/auth-proxy
RUN go build .
CMD ./auth-proxy
