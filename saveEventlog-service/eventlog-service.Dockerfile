#base go image
FROM golang:1.19-alpine as builder


WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /savelog .


#build multi stage
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /savelog /app/savelog

ENTRYPOINT [ "/app/savelog" ]