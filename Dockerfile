FROM node

EXPOSE 80

COPY ./app /srv/app
WORKDIR /srv/app

CMD ["node", "server.js"]
