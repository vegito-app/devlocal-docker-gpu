REACT_APP_VERSION = $(VERSION)

FRONTEND_BUILD_DIR = $(CURDIR)/application/frontend/build

application-frontend-build: application/frontend/node_modules
	@cd $(CURDIR)/application/frontend && npm --loglevel=verbose run build
.PHONY: application-frontend-build

$(FRONTEND_BUILD_DIR): application-frontend-build

UI_JAVASCRIPT_SOURCE_FILE = $(CURDIR)/application/frontend/build/bundle.js

application-frontend-bundle: 
	@cd $(CURDIR)/application/frontend && npm run dev:server
.PHONY: application-frontend-bundle

$(UI_JAVASCRIPT_SOURCE_FILE): application-frontend-bundle

application-frontend-start:
	@cd $(CURDIR)/application/frontend && npm start
.PHONY: application-frontend-start

application-frontend-npm-ci:
	@cd $(CURDIR)/application/frontend && npm ci
.PHONY: application-frontend-npm-ci