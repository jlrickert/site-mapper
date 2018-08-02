SOURCES = crawler.go graph.go main.go resources.go siteMap.go url.go util.go
RESOURCES = $(wildcard resources/*)

all: build

run: build
	./site-mapper $(filter-out $@, $(MAKECMDGOALS))

build: site-mapper

site-mapper: $(SOURCES)
	godep go build

resources.go: scripts/genResources.go $(RESOURCES)
	godep go generate

%:
	@true
