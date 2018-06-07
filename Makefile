SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

bin/smtp-client: $(SOURCES)
	go build -o bin/smtp-client
