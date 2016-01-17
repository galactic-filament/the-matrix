FROM node

RUN apt-get update -q \
  && apt-get install -yq netcat

COPY ./app /srv/app
WORKDIR /srv/app

RUN npm install -g --silent mocha tsd \
  && tsd -V \
  && mocha -V \
  && npm install --silent \
  && tsd install

CMD ["./bin/run-tests"]
