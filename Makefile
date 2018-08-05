BIN = site-mapper
SOURCES = $(wildcard main.go crawler/*.go)
RESOURCES = $(wildcard crawler/resources/*)

all: install

run: build
	./site-mapper $(filter-out $@, $(MAKECMDGOALS))

install:
	godep go install

build: site-mapper

site-mapper: $(SOURCES)
	godep go build -o $(BIN)

generate:
	godep go generate
	godep go generate github.com/jlrickert/site-mapper/crawler

%:
	@true
