FROM node

COPY ./app /srv/app
WORKDIR /srv/app

RUN npm install -g --silent mocha coffee-script \
  && npm install --silent

CMD ["npm", "run", "test"]
