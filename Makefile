.PHONY: default
default: build

all: build deploy

build:
	@nerdctl compose build

deploy:
	nerdctl compose up -d
