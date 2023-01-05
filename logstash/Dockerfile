FROM docker.elastic.co/logstash/logstash:7.10.2

RUN logstash-plugin install logstash-integration-rabbitmq
RUN logstash-plugin install logstash-output-elasticsearch

COPY logstash.conf /usr/share/logstash/pipeline/logstash.conf
