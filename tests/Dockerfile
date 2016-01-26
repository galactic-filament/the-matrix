FROM node

RUN apt-get update -q \
  && apt-get install -yq netcat

COPY ./app /srv/app
WORKDIR /srv/app

RUN npm install -g --silent typescript ts-node tsd \
  && tsc -v \
  && ts-node -v \
  && tsd -V \
  && npm install --silent \
  && tsd install

CMD ["./bin/run-tests"]
