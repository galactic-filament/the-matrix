FROM node

# installing docker and docker-compose
ENV DOCKER_COMPOSE_VERSION 1.4.2
RUN curl -sSL https://get.docker.com/ubuntu/ | sh \
  && curl -L https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose \
  && chmod +x /usr/local/bin/docker-compose

# adding github to known_hosts
ENV KNOWN_HOSTS_PATH /root/.ssh/known_hosts
RUN mkdir /root/.ssh \
  && touch $KNOWN_HOSTS_PATH \
  && ssh-keyscan -H github.com >> $KNOWN_HOSTS_PATH \
  && chmod 600 $KNOWN_HOSTS_PATH

# copying our code over
COPY ./app /srv/app
WORKDIR /srv/app

# installing nodejs assets
RUN npm install -g mocha \
  && npm install

CMD ["npm", "run", "test"]
