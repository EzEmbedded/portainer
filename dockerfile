
FROM golang:1.16.6-alpine AS development

ENV PROJECT_PATH=/portainer
ENV PATH=$PATH:$PROJECT_PATH/dist
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apk add --no-cache ca-certificates bash alpine-sdk nodejs npm yarn curl
RUN apk add --no-cache automake nasm autoconf build-base zlib zlib-dev libpng libpng-dev libwebp libwebp-dev libjpeg-turbo libjpeg-turbo-dev

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH
RUN go version && node -v && yarn -v
# RUN yarn add cypress --dev

RUN yarn
RUN  yarn build

FROM alpine:latest AS production
# FROM portainer/base AS production
# FROM alpine:3.13.2 AS production
# RUN apk --no-cache add ca-certificates
COPY --from=development /portainer/dist /
# USER nobody:nogroup
VOLUME /data
WORKDIR /
EXPOSE 9000
EXPOSE 8000
ENTRYPOINT ["/portainer"]