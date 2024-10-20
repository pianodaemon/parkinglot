ECHO = /bin/echo

BASE		:= $(SIDECARS_HOME)
OPERATIONAL	:= $(BASE)/DOS
REGISTRY	:= localhost:5000

# Inhouse images for local registry
PRICES_IMG        := $(REGISTRY)/sidecars/prices:latest

all:	prices_push

prices_push:	prices_build
	(docker --debug push $(PRICES_IMG))

prices_build:
	(cd $(OPERATIONAL)/prices && \
	docker --debug image build -f Dockerfile . \
	-t "$(PRICES_IMG)")

clean:
	- (docker rmi $(PRICES_IMG))
