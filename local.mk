
LATEST_BUILDER_IMAGE ?= $(PUBLIC_IMAGES_BASE):builder-latest

LOCAL_DIR ?= $(CURDIR)

LOCAL_DOCKER_COMPOSE ?= docker compose -f $(LOCAL_DIR)/docker-compose.yml

local-application-install: application-frontend-build application-frontend-bundle backend-install 
.PHONY: local-application-install

local-application-backend-install: $(APPLICATION_BACKEND_INSTALL_BIN) $(FRONTEND_BUILD_DIR) $(UI_JAVASCRIPT_SOURCE_FILE)
	@$(APPLICATION_BACKEND_INSTALL_BIN)
.PHONY: local-application-backend-install

BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE=$(LOCAL_DIR)/.containers/docker-buildx-cache/local-builder
$(BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE):;	@mkdir -p "$@"
ifneq ($(wildcard $(BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)/index.json),)
BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ = type=local,src=$(BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)
endif
BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE= type=local,dest=$(BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)

local-docker-compose-dev-image-pull:
	@$(LOCAL_DOCKER_COMPOSE) pull dev
.PHONY: local-docker-compose-dev-image-pull

local-docker-compose-dev-logs:
	@$(LOCAL_DOCKER_COMPOSE) logs dev
.PHONY: local-docker-compose-dev-logs

local-docker-compose-dev-logs-f:
	@$(LOCAL_DOCKER_COMPOSE) logs -f dev
.PHONY: local-docker-compose-dev-logs-f

LOCAL_DOCKER_COMPOSE_SERVICES ?= \
  android-studio \
  vault-dev \
  firebase-emulators \
  clarinet-devnet \
  application-backend \
  application-tests

local-docker-compose-up: $(LOCAL_DOCKER_COMPOSE_SERVICES)
.PHONY: local-docker-compose-up

local-docker-compose-rm-all: $(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-rm)
.PHONY: local-docker-compose-rm-all

$(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-image-pull):
	$(LOCAL_DOCKER_COMPOSE) pull $(@:local-%-image-pull=%)
.PHONY: $(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-image-pull)

$(LOCAL_DOCKER_COMPOSE_SERVICES):
	@$(MAKE) $(@:%=local-%-docker-compose-up)
.PHONY: $(LOCAL_DOCKER_COMPOSE_SERVICES)

$(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-rm): 
	@$(MAKE) $(@:%-rm=%-stop)
	@$(LOCAL_DOCKER_COMPOSE) rm -f $(@:local-%-docker-compose-rm=%)
.PHONY: $(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-rm)

$(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-start):
	@-$(LOCAL_DOCKER_COMPOSE) start $(@:local-%-docker-compose-start=%) 2>/dev/null
.PHONY: $(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-start)

$(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-stop):
	@-$(LOCAL_DOCKER_COMPOSE) stop $(@:local-%-docker-compose-stop=%) 2>/dev/null
.PHONY: $(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-stop)

$(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-logs):
	@$(LOCAL_DOCKER_COMPOSE) logs $(@:local-%-docker-compose-logs=%)
.PHONY: $(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-logs)

$(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-logs-f):
	@$(LOCAL_DOCKER_COMPOSE) logs --follow $(@:local-%-docker-compose-logs-f=%)
.PHONY: $(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-logs-f)

$(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-sh):
	@$(LOCAL_DOCKER_COMPOSE) exec -it $(@:local-%-docker-compose-sh=%) bash
.PHONY: $(LOCAL_DOCKER_COMPOSE_SERVICES:%=local-%-docker-compose-sh)

-include $(LOCAL_DIR)/docker/docker.mk
-include $(LOCAL_DIR)/android-studio/android-studio.mk
-include $(LOCAL_DIR)/clarinet-devnet/clarinet-devnet.mk
-include $(LOCAL_DIR)/github/github.mk
-include $(LOCAL_DIR)/firebase-emulators/firebase-emulators.mk
-include $(LOCAL_DIR)/vault-dev/vault-dev.mk
-include $(LOCAL_DIR)/application-tests/application-tests.mk
