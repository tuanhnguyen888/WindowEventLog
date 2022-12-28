#base go image
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

RUN go build -v -o brokerApp ./cmd/api


#build multi stage
FROM alpine:latest

RUN mkdir /app

COPY --from=builder brokerApp /app

ENTRYPOINT [ "/app/brokerApp" ]