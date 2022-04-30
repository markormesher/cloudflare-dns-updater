FROM node:16.14.2-alpine

WORKDIR /cloudflare-dns-updater

COPY ./package.json ./yarn.lock ./
RUN yarn install && yarn cache clean

COPY ./tsconfig.json ./
COPY ./src ./src
RUN yarn tsc

CMD node ./build/index.js
