ROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o listenerApp ./cmd/api

run chmod +x /app/listenerApp

#build multi stage
FROM alpine:latest

RUN mkdir /app

COPY --from=builder listenerApp /app

CMD [ "/app/listenerApp" ]