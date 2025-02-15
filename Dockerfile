FROM node:22.14.0-alpine@sha256:f96abbefa5558bcf3a309a5da9e1e9bd6dece7704b895c1213ca62f245061d3f AS builder

WORKDIR /app

COPY .yarn/ .yarn/
COPY package.json yarn.lock .yarnrc.yml .pnp* ./
RUN yarn install

COPY ./src ./src
COPY ./tsconfig.json ./

RUN yarn build

# ---

FROM node:22.14.0-alpine@sha256:f96abbefa5558bcf3a309a5da9e1e9bd6dece7704b895c1213ca62f245061d3f

WORKDIR /app

COPY .yarn/ .yarn/
COPY package.json yarn.lock .yarnrc.yml .pnp* ./
RUN yarn workspaces focus --all --production

COPY --from=builder /app/build /app/build

CMD yarn start
