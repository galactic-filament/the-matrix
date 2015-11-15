FROM node

RUN apt-get update -q \
  && apt-get install -yq netcat

COPY ./app /srv/app
WORKDIR /srv/app

RUN npm install -g --silent mocha coffee-script \
  && npm install --silent

CMD ["./bin/run-tests"]
