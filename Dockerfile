FROM node:20.3.0-alpine

WORKDIR /cloudflare-dns-updater

COPY ./package.json ./yarn.lock ./
RUN yarn install && yarn cache clean

COPY ./tsconfig.json ./
COPY ./src ./src
RUN yarn build

CMD yarn start
