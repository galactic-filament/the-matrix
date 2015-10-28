FROM node

EXPOSE 80

COPY ./app /srv/app
WORKDIR /srv/app

RUN npm install

CMD ["node", "server.js"]
