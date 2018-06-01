ORGANIZATION=jberlinsky
PROJECT_NAME=faxman-server
GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
RELEASE_DIR=bin
DEVTOOL_DIR=devtools
PACKAGE=github.com/$ORGANIZATION/$PROJECT_NAME
REVISION=$(shell git rev-parse --verify HEAD)
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./_vendor-*/")

.PHONY: clean build build-linux-amd64 build-linux-386 build-darwin-amd64 $(RELEASE_DIR)/$(PROJECT_NAME)_$(GOOS)_$(GOARCH) all

all: installdeps clean fmt simplify check build-linux-amd64 build-linux-386 build-darwin-amd64

build: $(RELEASE_DIR)/$(PROJECT_NAME)_$(GOOS)_$(GOARCH) $(RELEASE_DIR)/$(PROJECT_NAME)_worker_$(GOOS)_$(GOARCH)

build-linux-amd64:
	$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-386:
	@$(MAKE) build GOOS=linux GOARCH=386

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

$(RELEASE_DIR)/$(PROJECT_NAME)_$(GOOS)_$(GOARCH):
ifndef VERSION
	@echo '[ERROR] $$VERSION must be specified'
	exit 255
endif
	go build -ldflags "-X $(PACKAGE).rev=$(REVISION) -X $(PACKAGE).ver=$(VERSION)" \
		-o $(RELEASE_DIR)/$(PROJECT_NAME)_$(GOOS)_$(GOARCH)_$(VERSION) main.go

$(RELEASE_DIR)/$(PROJECT_NAME)_worker_$(GOOS)_$(GOARCH):
ifndef VERSION
	@echo '[ERROR] $$VERSION must be specified'
	exit 255
endif
	go build -ldflags "-X $(PACKAGE).rev=$(REVISION) -X $(PACKAGE).ver=$(VERSION)" \
		-o $(RELEASE_DIR)/$(PROJECT_NAME)_worker_$(GOOS)_$(GOARCH)_$(VERSION) worker.go

installdeps:
	go get -u github.com/golang/dep/cmd/dep
	@PATH=$(DEVTOOL_DIR)/$(GOOS)/$(GOARCH):$(PATH) dep ensure

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	# @for d in $$(go list ./... | grep -v /vendor/ | grep -v _vendor); do golint $${d}; done
	# @go tool vet ${SRC}

clean:
	rm -rf $(RELEASE_DIR)/$(PROJECT_NAME)_*
