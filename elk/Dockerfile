FROM java:8

ENV ELASTICSEARCH_VERSION 2.2.0
ENV LOGSTASH_VERSION 2.2.2-1_all
ENV KIBANA_VERSION 4.4.1-linux-x64

RUN apt-get update -q \
  && apt-get install -yq curl logrotate

# installing elasticsearch
RUN curl -L -O -s https://download.elastic.co/elasticsearch/elasticsearch/elasticsearch-${ELASTICSEARCH_VERSION}.deb \
  && dpkg -i elasticsearch-${ELASTICSEARCH_VERSION}.deb \
  && rm ./elasticsearch-${ELASTICSEARCH_VERSION}.deb

# installing logstash
RUN curl -L -O -s https://download.elastic.co/logstash/logstash/packages/debian/logstash_${LOGSTASH_VERSION}.deb \
  && dpkg -i logstash_${LOGSTASH_VERSION}.deb \
  && rm ./logstash_${LOGSTASH_VERSION}.deb

# installing kibana
RUN curl -L -O -s https://download.elastic.co/kibana/kibana/kibana-${KIBANA_VERSION}.tar.gz \
  && mkdir /srv/kibana \
  && tar --strip-components=1 -C /srv/kibana -xf ./kibana-${KIBANA_VERSION}.tar.gz \
  && rm ./kibana-${KIBANA_VERSION}.tar.gz

CMD ["bash"]
