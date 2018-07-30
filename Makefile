SOURCES = crawler.go main.go resources.go siteMap.go url.go util.go

all: build

run: build
	./site-mapper $(filter-out $@, $(MAKECMDGOALS))

build: site-mapper

site-mapper: $(SOURCES)
	go build

resources.go: scripts/genResources.go resources/.*
	go generate

%:
	@true
