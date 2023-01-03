FROM docker.elastic.co/logstash/logstash:8.5.3

RUN rm -f /usr/share/logstash/pipelines/logstash.conf

COPY pipelines/ /usr/share/logstash/pipelines/

COPY config/ /usr/share/logstash/config/