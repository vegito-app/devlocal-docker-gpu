GITHUB_ACTIONS_RUNNER_STACK_ID = $(shell echo $$RANDOM)
GITHUB_ACTIONS_RUNNER_STACK = github-actions-$(GITHUB_ACTIONS_RUNNER_STACK_ID)

GITHUB_ACTIONS_RUNNER_IMAGE = $(PUBLIC_IMAGES_BASE):github-action-runner-latest
GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE=$(CURDIR)/.containers/docker-buildx-cache/infra-github
$(GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE):;	@mkdir -p "$@"
ifneq ($(wildcard $(GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)/index.json),)
GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ = type=local,src=$(GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)
endif
GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE= type=local,dest=$(GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)

# Build image for local run. This target will not push an image to the distant registry.
github-action-runner-image: $(GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE) docker-buildx-setup
	@$(DOCKER_BUILDX_BAKE) --print github-action-runner
	@$(DOCKER_BUILDX_BAKE) --load github-action-runner
.PHONY: github-action-runner-image

# Build image for local run and push it.
github-action-runner-image-push: $(GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE) docker-buildx-setup
	@$(DOCKER_BUILDX_BAKE) --print github-action-runner
	@$(DOCKER_BUILDX_BAKE) --push github-action-runner
.PHONY: github-action-runner-image-push

# This target will build and push a multi architecture image.
github-action-runner-image-ci: docker-buildx-setup
	@$(DOCKER_BUILDX_BAKE) --print github-action-runner-ci
	@$(DOCKER_BUILDX_BAKE) --push github-action-runner-ci
.PHONY: github-action-runner-image-ci

GITHUB_DOCKER_COMPOSE := COMPOSE_PROJECT_NAME=$(GOOGLE_CLOUD_PROJECT_ID)-github-actions \
  docker compose -f $(CURDIR)/github/docker-compose.yml

github-action-runner-token-exist:
	@if [ ! -v GITHUB_ACTIONS_RUNNER_TOKEN ] ; then \
		echo missing GITHUB_ACTIONS_RUNNER_TOKEN environment variable. \
		exit -1 ; \
	fi
.PHONY: github-action-runner-token-exist

local-github-action-runner-docker-compose-up: github-action-runner-docker-compose-rm github-action-runner-token-exist
	@$(GITHUB_DOCKER_COMPOSE) up -d github-action-runner
.PHONY: github-action-runner-docker-compose-up

# This avoids github dangling offline containers on github.com side.
# It uses './config.sh remove' from github-action-runner containers entrypoint
# before exit (see ./entrypoint.sh)
local-github-action-runner-docker-compose-rm:
	@-$(GITHUB_DOCKER_COMPOSE) rm -s -f github-action-runner
.PHONY: local-github-action-runner-docker-compose-rm

LOCAL_GITHUB_WORKFLOWS_DIR := $(CURDIR)/.github/workflows/
LOCAL_GITHUB_ACT_SECRET_FILE := $(LOCAL_GITHUB_WORKFLOWS_DIR)/.secret

$(LOCAL_GITHUB_ACT_SECRET_FILE):
	@-rm -f $(LOCAL_GITHUB_ACT_SECRET_FILE) 2>&1
	@-echo DEV_GCLOUD_SERVICE_KEY=$$(jq -c . $(CURDIR)/infra/environments/gcloud-credentials.json) >> $(LOCAL_GITHUB_ACT_SECRET_FILE)
	@-echo STAGING_GCLOUD_SERVICE_KEY=$$(jq -c . $(CURDIR)/infra/environments/staging/gcloud-credentials.json) >> $(LOCAL_GITHUB_ACT_SECRET_FILE) 2>/dev/null 
	@-echo PRODUCTION_GCLOUD_SERVICE_KEY=$$(jq -c . $(CURDIR)/infra/environments/prod/gcloud-credentials.json) >> $(LOCAL_GITHUB_ACT_SECRET_FILE) 2>/dev/null

GITHUB_WORKFLOWS := \
  dev.yml \
  dev-feature.yml \
  staging.yml \
  production.yml

LOCAL_GITHUB_ACT := act --secret-file $(LOCAL_GITHUB_ACT_SECRET_FILE) \
	 -P self-hosted=$(GITHUB_ACTIONS_RUNNER_IMAGE) \
	 
$(GITHUB_WORKFLOWS:%=local-github-run-%-workflow): $(LOCAL_GITHUB_ACT_SECRET_FILE)
	@$(LOCAL_GITHUB_ACT) -W $(CURDIR)/.github/workflows/$(@:github-run-%-workflow=%)
.PHONY: $(GITHUB_WORKFLOWS:%=local-github-run-%-workflow)