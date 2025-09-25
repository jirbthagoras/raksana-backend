FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./out .

FROM alpine:latest

WORKDIR /main

COPY --from=builder /build/out ./out

RUN chmod +x ./out

EXPOSE 3000

CMD ["./out"]
