FROM node:21.7.3-alpine@sha256:0a50081b5723b3cfe2ef3a3c5675906b0bb942a4b8ede1f6ba5be6ec88413ec4 AS builder

WORKDIR /app

COPY .yarn/ .yarn/
COPY package.json yarn.lock .yarnrc.yml .pnp* ./
RUN yarn install

COPY ./src ./src
COPY ./tsconfig.json ./

RUN yarn build

# ---

FROM node:21.7.3-alpine@sha256:0a50081b5723b3cfe2ef3a3c5675906b0bb942a4b8ede1f6ba5be6ec88413ec4

WORKDIR /app

COPY .yarn/ .yarn/
COPY package.json yarn.lock .yarnrc.yml .pnp* ./
RUN yarn workspaces focus --all --production

COPY --from=builder /app/build /app/build

CMD yarn start
