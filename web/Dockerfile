FROM node:lts-alpine AS builder

WORKDIR /build

COPY package.json yarn.lock /build/
RUN yarn install

COPY . /build
RUN yarn build && yarn build-html


FROM nginx:stable-alpine

COPY --from=builder /build/dist/ /usr/share/nginx/html/
COPY default.conf /etc/nginx/conf.d/
