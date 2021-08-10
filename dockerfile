# FROM ubuntu

# # Expose port for the Portainer UI and Edge server
# EXPOSE 9000
# EXPOSE 8000

# WORKDIR /src/portainer

FROM golang:1.16.6-alpine AS development

ENV PROJECT_PATH=/portainer
ENV PATH=$PATH:$PROJECT_PATH/dist
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

# Set TERM as noninteractive to suppress debconf errors
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections

RUN apk add --no-cache ca-certificates bash alpine-sdk

# # Set default go version
# ARG GO_VERSION=go1.16.6.linux-amd64

# Install packages
RUN apt-get update --fix-missing && apt-get install -qq \
	dialog \
	apt-utils \
	curl \
	build-essential \
	nodejs \
	git \
	wget

# Install Yarn
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
	&& echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list \
	&& apt-get update && apt-get -y install yarn

# # Install Golang
# RUN cd /tmp \
# 	&& wget -q https://dl.google.com/go/${GO_VERSION}.tar.gz \
# 	&& tar -xf ${GO_VERSION}.tar.gz \
# 	&& mv go /usr/local

# # Configure Go
# ENV PATH "$PATH:/usr/local/go/bin"

# Confirm installation
# RUN node version && yarn -v


RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH
RUN  build

FROM alpine:3.13.2 AS production
RUN apk --no-cache add ca-certificates
COPY --from=development /portainer/dist /portainer
USER nobody:nogroup
EXPOSE 9000
EXPOSE 8000
ENTRYPOINT ["/portainer"]