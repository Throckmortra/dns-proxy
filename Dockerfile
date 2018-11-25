FROM golang:latest as builder
RUN go get github.com/Masterminds/glide
WORKDIR /go/src/app
ADD glide.yaml glide.yaml
ADD glide.lock glide.lock
RUN glide install
ADD . src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main src/main.go

FROM alpine:3.7
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root
COPY --from=builder /go/src/app/main .
ADD config.toml config.toml
ADD server.rsa.crt server.rsa.crt
ADD server.rsa.key server.rsa.key
CMD ["./main"]
EXPOSE 53
