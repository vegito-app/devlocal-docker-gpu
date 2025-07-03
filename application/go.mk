GO_VERSION = 1.24

GO_MODULES := \
	application/backend \
	application/images \
	application/firebase/functions/auth \
	proxy

go-mod-tidy: $(GO_MODULES:%=go-%-mod-tidy)
.PHONY: go-mod-tidy

$(GO_MODULES:%=go-%-mod-tidy):
	@cd $(CURDIR)/$(@:go-%-mod-tidy=%) && go mod tidy -v -go=$(GO_VERSION)
.PHONY: $(GO_MODULES:%=go-%-mod-tidy) 

go-mod-download: $(GO_MODULES:%=go-%-mod-download)
.PHONY: go-mod-download

$(GO_MODULES:%=go-%-mod-download):
	@cd $(CURDIR)/$(@:go-%-mod-download=%) && go mod download
.PHONY: $(GO_MODULES:%=go-%-mod-download) 

go-mod-upgrade: $(GO_MODULES:%=go-%-mod-upgrade)
.PHONY: go-mod-upgrade

$(GO_MODULES:%=go-%-mod-upgrade):
	@cd $(CURDIR)/$(@:go-%-mod-upgrade=%) && rm -rf vendor && go get -u -v ./...
.PHONY: $(GO_MODULES:%=go-%-mod-upgrade)

go-mod-vendor: $(GO_MODULES:%=go-%-mod-vendor)
.PHONY: go-mod-vendor

go-mod-vendor-rm: $(GO_MODULES:%=go-%-mod-vendor-rm)
.PHONY: go-mod-vendor-rm

$(GO_MODULES:%=go-%-mod-vendor):
	@cd $(CURDIR)/$(@:go-%-mod-vendor=%) && go mod vendor -v
.PHONY: $(GO_MODULES:%=go-%-mod-vendor) 

$(GO_MODULES:%=go-%-mod-vendor-rm):
	@rm -rf $(CURDIR)/$(@:go-%-mod-vendor-rm=%)/vendor
.PHONY: $(GO_MODULES:%=go-%-mod-vendor-rm) 
