FROM golang

# # installing docker and docker-compose
# ENV DOCKER_COMPOSE_VERSION 1.4.2
# RUN curl -sSL https://get.docker.com/ | sh \
#   && curl -L https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose \
#   && chmod +x /usr/local/bin/docker-compose

# adding github to known_hosts
ENV KNOWN_HOSTS_PATH /root/.ssh/known_hosts
RUN mkdir /root/.ssh \
  && touch $KNOWN_HOSTS_PATH \
  && ssh-keyscan -H github.com >> $KNOWN_HOSTS_PATH \
  && chmod 600 $KNOWN_HOSTS_PATH

# copying our code over
ENV APP_PATH github.com/ihsw/the-matrix/app
ADD ./app ./src/$APP_PATH
RUN go get ./src/$APP_PATH/... \
  && go get -t $APP_PATH \
  && go install $APP_PATH

CMD ["./bin/app"]
