FROM golang:alpine AS builder

RUN go version

WORKDIR /temp

COPY . .
RUN go mod download

RUN GOOS=linux go build -o ./.bin/rusprofile ./cmd/rusprofile/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder ./temp/.bin/rusprofile .
COPY --from=builder ./temp/config config/
CMD ["./rusprofile"]