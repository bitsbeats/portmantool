FROM node:lts-alpine AS builder

COPY . /build

WORKDIR /build
RUN yarn install && yarn build && yarn build-html


FROM nginx:stable-alpine

COPY --from=builder /build/dist/ /usr/share/nginx/html/
COPY default.conf /etc/nginx/conf.d/
