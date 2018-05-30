# The import path is where your repository can be found.
# To import subpackages, always prepend the full import path.
# If you change this, run `make clean`. Read more: https://git.io/vM7zV
IMPORT_PATH := sql-client

V := 1 # When V is set, print commands and build progress.

.PHONY: all
all: sql-client

.PHONY: sql-client
sql-client: .GOPATH/.ok
	$Q go install -tags netgo $(IMPORT_PATH)/cmd
	$Q mv bin/cmd bin/sql-client
	
.PHONY: update
update: .GOPATH/.ok
	$Q glide mirror set https://golang.org/x/sys https://github.com/golang/sys
	$Q glide mirror set https://golang.org/x/net https://github.com/golang/net
	$Q glide up -v

.PHONY: clean

clean:
	$Q rm -rf bin/* .GOPATH

export GOPATH := $(CURDIR)/.GOPATH

unexport GOBIN

Q := $(if $V,,@)

.GOPATH/.ok:
	$Q rm -rf $(CURDIR)/.GOPATH
	$Q mkdir -p $(CURDIR)/.GOPATH/src
	$Q ln -sf $(CURDIR) $(CURDIR)/.GOPATH/src/$(IMPORT_PATH)
	$Q mkdir -p $(CURDIR)/bin
	$Q ln -sf $(CURDIR)/bin $(CURDIR)/.GOPATH/bin
	$Q touch $@