GOOGLE_CLOUD_REGION ?= europe-west1
REGISTRY ?= $(GOOGLE_CLOUD_REGION)-docker.pkg.dev
PUBLIC_REPOSITORY ?= $(REGISTRY)/$(GOOGLE_CLOUD_PROJECT_ID)/docker-repository-public
PUBLIC_IMAGES_BASE ?= $(PUBLIC_REPOSITORY)/moov-dev-439608

DOCKER_BUILDX_BAKE_APPLICATION_IMAGES_WORKERS_IMAGES = \
  application-images-cleaner  \
  application-images-moderator 

DOCKER_BUILDX_BAKE_APPLICATION_IMAGES = \
  application-backend

DOCKER_BUILDX_BAKE_LOCAL_IMAGES = \
  android-studio \
  clarinet-devnet \
  application-tests \
  firebase-emulators \
  vault-dev 

docker-local-images-pull: $(DOCKER_BUILDX_BAKE_LOCAL_IMAGES:%=local-%-image-pull) local-docker-compose-dev-image-pull
.PHONY: docker-local-images-pull

docker-local-images-push: $(DOCKER_BUILDX_BAKE_LOCAL_IMAGES:%=local-%-image-push) local-builder-image-push
.PHONY: docker-local-images-push

DOCKER_BUILDX_BAKE_IMAGES ?= \
  $(DOCKER_BUILDX_BAKE_APPLICATION_IMAGES) \
  $(DOCKER_BUILDX_BAKE_APPLICATION_IMAGES_WORKERS_IMAGES) \
  $(DOCKER_BUILDX_BAKE_LOCAL_IMAGES) 

DOCKER_BUILDX_BAKE ?= docker buildx bake \
	-f $(LOCAL_DIR)/docker/docker-bake.hcl \
	-f $(LOCAL_DIR)/docker-bake.hcl \
	$(DOCKER_BUILDX_BAKE_LOCAL_IMAGES:%=-f $(LOCAL_DIR)/%/docker-bake.hcl) \
	$(DOCKER_BUILDX_BAKE_APPLICATION_IMAGES_WORKERS_IMAGES:application-images-%=-f $(LOCAL_DIR)/application/images/%/docker-bake.hcl) \
	$(DOCKER_BUILDX_BAKE_APPLICATION_IMAGES:application-%=-f $(LOCAL_DIR)/application/%/docker-bake.hcl) \
	-f $(LOCAL_DIR)/github/docker-bake.hcl

docker-images-ci-multi-arch: docker-buildx-setup local-builder-image-ci
	@$(DOCKER_BUILDX_BAKE) --print services-push-multi-arch
	@$(DOCKER_BUILDX_BAKE) --push services-push-multi-arch
.PHONY: docker-images-ci-multi-arch

docker-images-local-arch: local-builder-image
	$(DOCKER_BUILDX_BAKE) --print services-load-local-arch
	$(DOCKER_BUILDX_BAKE) --load services-load-local-arch
.PHONY: docker-images-local-arch

docker-buildx-setup: 
	@-docker buildx create --name $(GOOGLE_CLOUD_PROJECT_ID)-builder 2>/dev/null 
	@-docker buildx use $(GOOGLE_CLOUD_PROJECT_ID)-builder 2>/dev/null 
.PHONY: docker-buildx-setup

docker-login: gcloud-auth-docker
	@docker login $(REGISTRY)/$(GOOGLE_CLOUD_PROJECT_ID)
.PHONY: docker-login

docker-sock:
	sudo chmod o+rw /var/run/docker.sock
.PHONY: docker-sock

$(DOCKER_BUILDX_BAKE_IMAGES:%=local-%-image): docker-buildx-setup
	@$(DOCKER_BUILDX_BAKE) --print $(@:local-%-image=%)
	@$(DOCKER_BUILDX_BAKE) --load $(@:local-%-image=%)
.PHONY: $(DOCKER_BUILDX_BAKE_IMAGES:%=local-%-image)

$(DOCKER_BUILDX_BAKE_IMAGES:%=local-%-image-push):
	@$(DOCKER_BUILDX_BAKE) --print $(@:local-%-image-push=%)
	@$(DOCKER_BUILDX_BAKE) --push $(@:local-%-image-push=%)
.PHONY: $(DOCKER_BUILDX_BAKE_IMAGES:%=local-%-image-push)

# $(DOCKER_BUILDX_BAKE_IMAGES:%=local-%-image-pull):
# 	@$(DOCKER_BUILDX_BAKE) --print $(@:local-%-image-pull=%)
# 	@$(DOCKER_BUILDX_BAKE) --pull $(@:local-%-image-pull=%)
# .PHONY: $(DOCKER_BUILDX_BAKE_IMAGES:%=local-%-image-pull)

$(DOCKER_BUILDX_BAKE_IMAGES:%=local-%-image-ci): docker-buildx-setup
	@$(DOCKER_BUILDX_BAKE) --print $(@:local-%-image-ci=%-ci)
	@$(DOCKER_BUILDX_BAKE) --push $(@:local-%-image-ci=%-ci)
.PHONY: $(DOCKER_BUILDX_BAKE_IMAGES:%=local-%-image-ci)

local-builder-image: $(BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE) docker-buildx-setup
	$(DOCKER_BUILDX_BAKE) --print builder
	$(DOCKER_BUILDX_BAKE) --load builder
.PHONY: local-builder-image

local-builder-image-push: docker-buildx-setup
	@$(DOCKER_BUILDX_BAKE) --print builder
	@$(DOCKER_BUILDX_BAKE) --push builder
.PHONY: local-builder-image-push

local-builder-image-ci: docker-buildx-setup
	@$(DOCKER_BUILDX_BAKE) --print builder-ci
	@$(DOCKER_BUILDX_BAKE) --push builder-ci
.PHONY: local-builder-image-ci
