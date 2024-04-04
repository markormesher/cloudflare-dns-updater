FROM node:21.7.2-alpine@sha256:c4679fe0dc0434d1c002d0a453a8b09a397d93e44c6c12dd1ec7affa3891422d AS builder

WORKDIR /app

COPY .yarn/ .yarn/
COPY package.json yarn.lock .yarnrc.yml .pnp* ./
RUN yarn install

COPY ./src ./src
COPY ./tsconfig.json ./

RUN yarn build

# ---

FROM node:21.7.2-alpine@sha256:c4679fe0dc0434d1c002d0a453a8b09a397d93e44c6c12dd1ec7affa3891422d

WORKDIR /app

COPY .yarn/ .yarn/
COPY package.json yarn.lock .yarnrc.yml .pnp* ./
RUN yarn workspaces focus --all --production

COPY --from=builder /app/build /app/build

CMD yarn start
