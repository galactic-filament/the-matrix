FROM node

COPY ./app /srv/app
WORKDIR /srv/app

RUN npm install -g mocha coffee-script \
  && npm install

CMD ["npm", "run", "test"]
