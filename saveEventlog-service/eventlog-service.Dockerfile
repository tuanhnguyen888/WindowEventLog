FROM alpine:latest

RUN mkdir /app

COPY saveLogApp /app

CMD [ "/app/saveLogApp"]