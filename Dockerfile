FROM golang:stretch as builder
WORKDIR /go/src/github.com/stanleynguyen/web3-lucky-draw
RUN apt update && apt upgrade -y
COPY . .
RUN go build -o app .

FROM debian:stretch
ENV GO_ENV=production
RUN apt update && apt upgrade -y
WORKDIR /root/
COPY --from=builder /go/src/github.com/stanleynguyen/web3-lucky-draw/app .
CMD ["./app"]
