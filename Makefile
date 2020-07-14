
include .env

#GIT VERSION INFO FOR PACKAGING
COMMITID= $(shell git rev-parse --short HEAD)
BUILDDATE=$(shell date -u +'%Y-%m-%dT%H-%M-%SZ')
REFID= $(shell git rev-parse --abbrev-ref HEAD)

# Go related variables.
GOMODULENAME="ovsdb-statsd-client"
GOBASE=$(shell pwd)
GOPATH="$(GOBASE):$(GOBASE)/vendor"
GOBIN=$(GOBASE)/bin
DESTDIR?=/usr/local
GOFILES=$(shell find "$(GOBASE)" -type f -name '*.go')
GOSRC=$(GOBASE)
GOLOG=$(GOBASE)/logs
GOPKG=$(GOBASE)/pkg

GOMAIN=$(GOSRC)/cmd/client-main
PROJECTNAME=$(shell basename "$(GOMAIN)")

# Redirect error output to a file, so we can show it in development mode.
STDERR=$(GOLOG)/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID=$(GOBIN)/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

## install-dep : Install missing dependencies. Runs `go get` internally. e.g; make install-dep get=github.com/foo/bar
install-dep:
	@-touch $(STDERR)
	@echo "  >  Detailed error logs are at $(STDERR)"
	@$(MAKE) go-get 2>&1 | tee -a $(STDERR)

## install : Install ovsdb-statsd client binary. Custom location is possible to use. e.g; make install DESTDIR=/home/sugesh/test
install:
	@-touch $(STDERR)
	@$(MAKE) go-install 2>&1 | tee -a $(STDERR)

## uninstall : uninstall the ovsdb-statsd client binary. Custom location is possible to use. e.g; make install DESTDIR=/home/sugesh/test
uninstall: go-uninstall

## compile : Compile the binary
compile:
	@-touch $(STDERR)
	@echo "  >  Detailed error logs are at $(STDERR)"
	@$(MAKE) go-compile 2>&1 | tee -a $(STDERR)

## compile-debug : compile the binary in debug mode
compile-debug:
	@-touch $(STDERR)
	GCFLAGS='-race -gcflags=all="-N -l"'
	@echo "  >  Detailed error logs are at $(STDERR)"
	@$(MAKE) go-debug-compile 2>&1 | tee -a $(STDERR)

## clean : Clean build files. Runs `go clean` internally
clean:
	@echo "  >  Detailed error logs are at $(STDERR)"
	@-$(MAKE) go-clean 2>&1 | tee -a $(STDERR)


go-compile: go-clean go-get go-build

go-debug-compile: go-clean go-get go-debug-build

go-build:
	@echo "  >  Building binary..."
	@GO111MODULE=on CGO_ENABLED=1 CC=$(CC) GOBIN=$(GOBIN) $(GO) build $(GCFLAGS) -v -x -mod=vendor -o $(GOBIN)/$(PROJECTNAME) $(GOMAIN)

go-debug-build:
	@echo "  >  Building binary in debug mode..."
	@GO111MODULE=on CGO_ENABLED=1 CC=$(CC) GOBIN=$(GOBIN) $(GO) build -race $(GCFLAGS) -v -x -mod=vendor -o $(GOBIN)/$(PROJECTNAME) $(GOMAIN)

go-x86-64-linux-build:
	@echo "  >  Building binary for linux on ARM64..."
	# All supported architecutres can be found by 'go tool dist list'
	@GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=1 CC=$(CC) GOBIN=$(GOBIN) $(GO) build $(GCFLAGS) -v -x -mod=vendor -o $(GOBIN)/$(PROJECTNAME) $(GOMAIN)

go-generate:
	@echo "  >  Generating dependency files..."
	@GO111MODULE=on GOBIN=$(GOBIN) $(GO) generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
ifndef get

#	@GO111MODULE=on GOBIN=$(GOBIN) $(GO) get -v -u golang.org/x/net/icmp
#	@GO111MODULE=on GOBIN=$(GOBIN) $(GO) get -v -u golang.org/x/net/ipv4
#	@GO111MODULE=on GOBIN=$(GOBIN) $(GO) get -v -u golang.org/x/net/ipv6
# Get all omported modules into vendor.
	@GO111MODULE=on $(GO) mod vendor -v
else
	@GO111MODULE=on GOBIN=$(GOBIN) $(GO) get $(get)
endif

go-install:
	@mkdir -pv $(DESTDIR)/bin
	@mkdir -pv $(DESTDIR)/logs
	@mkdir -pv $(DESTDIR)/config
	@echo "  >  Installing ovsdb-statsd clietn at $(DESTDIR)"
	@GOBIN=$(DESTDIR)/bin $(GO) install $(GCFLAGS) $(GOMAIN)
	@cp -rv $(GOBASE)/config/* $(DESTDIR)/config

go-uninstall:
	@echo "  >  Uninstalling ovsdb-statsd client from $(DESTDIR)"
	@-rm -rvf $(DESTDIR)/bin/*
	@-rm -rvf $(DESTDIR)/logs/*
	@-rm -rvf $(DESTDIR)/config/*

go-clean:
	@echo "  >  Cleaning build cache"
	@-GO111MODULE=on GOBIN=$(GOBIN) $(GO) clean -r -cache $(GOMAIN)
	@cd $(GOBIN) && rm -rf *
	@rm -rvf $(GOLOG)/*
	@rm -rf $(TMPGODOCDIR)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo