<p align="center">
  <img title="portainer" src='https://github.com/portainer/portainer/blob/develop/app/assets/images/portainer-github-banner.png?raw=true' />
</p>

**Portainer CE** is a lightweight ‘universal’ management GUI that can be used to **easily** manage Docker, Swarm, Kubernetes and ACI environments. It is designed to be as **simple** to deploy as it is to use.

Portainer consists of a single container that can run on any cluster. It can be deployed as a Linux container or a Windows native container.

**Portainer** allows you to manage all your orchestrator resources (containers, images, volumes, networks and more) through a super-simple graphical interface.

A fully supported version of Portainer is available for business use. Visit http://www.portainer.io to learn more

## Demo

You can try out the public demo instance: http://demo.portainer.io/ (login with the username **admin** and the password **tryportainer**).

Please note that the public demo cluster is **reset every 15min**.

## Latest Version

Portainer CE is updated regularly. We aim to do an update release every couple of months.

**The latest version of Portainer is 2.6.x** And you can find the release notes [here.](https://www.portainer.io/blog/new-portainer-ce-2.6.0-release)
Portainer is on version 2, the second number denotes the month of release.

## Getting started

- [Deploy Portainer](https://documentation.portainer.io/quickstart/)
- [Documentation](https://documentation.portainer.io)
- [Contribute to the project](https://documentation.portainer.io/contributing/instructions/)

## Features & Functions

View [this](https://www.portainer.io/products) table to see all of the Portainer CE functionality and compare to Portainer Business.

- [Portainer CE for Docker / Docker Swarm](https://www.portainer.io/solutions/docker)
- [Portainer CE for Kubernetes](https://www.portainer.io/solutions/kubernetes-ui)
- [Portainer CE for Azure ACI](https://www.portainer.io/solutions/serverless-containers)

## Getting help

Portainer CE is an open source project and is supported by the community. You can buy a supported version of Portainer at portainer.io

Learn more about Portainers community support channels [here.](https://www.portainer.io/help_about)

- Issues: https://github.com/portainer/portainer/issues
- Slack (chat): https://portainer.io/slack/

You can join the Portainer Community by visiting community.portainer.io. This will give you advance notice of events, content and other related Portainer content.

## Reporting bugs and contributing

- Want to report a bug or request a feature? Please open [an issue](https://github.com/portainer/portainer/issues/new).
- Want to help us build **_portainer_**? Follow our [contribution guidelines](https://documentation.portainer.io/contributing/instructions/) to build it locally and make a pull request. We need all the help we can get!

## Security

- Here at Portainer, we believe in [responsible disclosure](https://en.wikipedia.org/wiki/Responsible_disclosure) of security issues. If you have found a security issue, please report it to <security@portainer.io>.

## Privacy

**To make sure we focus our development effort in the right places we need to know which features get used most often. To give us this information we use [Matomo Analytics](https://matomo.org/), which is hosted in Germany and is fully GDPR compliant.**

When Portainer first starts, you are given the option to DISABLE analytics. If you **don't** choose to disable it, we collect anonymous usage as per [our privacy policy](https://www.portainer.io/documentation/in-app-analytics-and-privacy-policy/). **Please note**, there is no personally identifiable information sent or stored at any time and we only use the data to help us improve Portainer.

## Limitations

Portainer supports "Current - 2 docker versions only. Prior versions may operate, however these are not supported.

## Licensing

Portainer is licensed under the zlib license. See [LICENSE](./LICENSE) for reference.

Portainer also contains code from open source projects. See [ATTRIBUTIONS.md](./ATTRIBUTIONS.md) for a list.

## 20210820

docker run hello-world
docker login
docker tag hello-world:latest scanlidocker/loraserver-app:v0.21
docker push scanlidocker/loraserver-app:v0.21

docker tag scanlidocker/chirpstack-application-server chirpstack-application-server

docker tag scanlidocker/chirpstack-network-server:v0.21 chirpstack-network-server

docker tag scanlidocker/chirpstack-gateway-bridge:v0.21 chirpstack-gateway-bridge
df -hl
apt install docker.io -y
docker run -d -p 1880:1880 --restart=always -v node_red_data:/data --name lora_nodered nodered/node-red
docker run -d -p 9000:9000 --restart=always -v /var/run/docker.sock:/var/run/docker.sock --name lora_portainer scanlidocker/portainer:v0.1

dd if=/dev/zero of=/dev/mmcblk2p2 bs=1024 seek=8 count=1

nmcli connection modify 'Wired connection 1' connection.autoconnect yes ipv4.method manual ipv4.address 192.168.14.14 ipv4.gateway 192.168.14.1 ipv4.dns 114.114.114.114

## portainer 开发环境搭建：

    1、在项目根目录下添加 dockerfile
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

    2、在.github/workflows下，添加buildimage_push.yml:
    	name: CI
    	on:
    	  workflow_dispatch:

    	jobs:
    	  buildimage:
    		runs-on: ubuntu-latest
    		steps:
    		  -
    			name: Checkout
    			uses: actions/checkout@v2
    			# with:
    			#   ref: 1.24
    		  -
    			name: Docker meta
    			id: meta
    			uses: docker/metadata-action@v3
    			with:
    			  images: |
    				scanlidocker/portainer
    			  tags: |
    				type=semver,pattern={{version}}
    				type=semver,pattern={{major}}
    				type=semver,pattern={{major}}.{{minor}}
    		  -
    			name: Set up QEMU
    			uses: docker/setup-qemu-action@v1
    		  -
    			name: Set up Docker Buildx
    			uses: docker/setup-buildx-action@v1
    		  -
    			name: Login to DockerHub
    			uses: docker/login-action@v1
    			with:
    			  username: ${{ secrets.DOCKER_HUB_USER_NAME }}
    			  password: ${{ secrets.DOCKER_HUB_PASSWORD }}
    		  -
    			name: Build and push
    			id: docker_build
    			uses: docker/build-push-action@v2
    			with:
    			  # platforms: linux/amd64
    			  platforms: linux/arm/v7
    			  # platforms: linux/arm64/v8
    			  context: .
    			  push: true
    			  tags: scanlidocker/portainer:v0.1
    			  labels: ${{ steps.meta.outputs.labels }}
    		  -
    			name: Image digest
    			run: echo ${{ steps.docker_build.outputs.digest }}
    3、在portainer\build\build_binary.sh里添加：
    			#!/bin/sh
    4、在portainer\package.json里去掉cypress相关的软件包依赖。
    5、core.js下载有时不成功，重新运行即可。
    6、架构支持：linux/arm/v7 linux/arm64/v8 linux/amd64
