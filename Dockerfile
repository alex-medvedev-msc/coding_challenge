FROM golang:1.10-alpine

RUN mkdir -p /go/src/github.com/messwith/coding_challenge

WORKDIR /go/src/github.com/messwith/coding_challenge

COPY . .

RUN go build -o run cmd/run.go

CMD ["./run"]

