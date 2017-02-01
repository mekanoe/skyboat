# TODO: don't use latest
FROM quay.io/saving-light/base:latest

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 7635B973 && \
	echo "deb http://ppa.launchpad.net/ubuntu-lxc/lxd-stable/ubuntu xenial main" | tee -a /etc/apt/sources.list && \
	apt-get update && \
	apt-get install -y --no-install-recommends \
		golang