variable "BUILDER_IMAGE_VERSION" {
  default = notequal("dev", VERSION) ? "${PUBLIC_IMAGES_BASE}:builder-${VERSION}" : ""
}

variable "LATEST_BUILDER_IMAGE" {
  default = "${PUBLIC_IMAGES_BASE}:builder-latest"
}

target "builder-ci" {
  args = {
    docker_version         = DOCKER_VERSION
    docker_compose_version = DOCKER_COMPOSE_VERSION
    docker_buildx_version  = DOCKER_BUILDX_VERSION
    terraform_version      = TERRAFORM_VERSION
    kubectl_version        = KUBECTL_VERSION
    node_version           = NODE_VERSION
    nvm_version            = NVM_VERSION
  }
  dockerfile = "Dockerfile"
  tags = [
    LATEST_BUILDER_IMAGE,
    notequal("", VERSION) ? BUILDER_IMAGE_VERSION : "",
  ]
  cache-from = [LATEST_BUILDER_IMAGE]
  cache-to   = ["type=inline"]
  platforms  = platforms
}

variable "BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE" {
  description = "local write cache for builder image build"
}

variable "BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ" {
  description = "local read cache for builder image build (cannot be used before first write)"
}

target "builder" {
  args = {
    docker_version         = DOCKER_VERSION
    docker_compose_version = DOCKER_COMPOSE_VERSION
    docker_buildx_version  = DOCKER_BUILDX_VERSION
    terraform_version      = TERRAFORM_VERSION
    kubectl_version        = KUBECTL_VERSION
    node_version           = NODE_VERSION
    nvm_version            = NVM_VERSION
  }
  dockerfile = "Dockerfile"
  tags = [
    LATEST_BUILDER_IMAGE,
    notequal("", VERSION) ? BUILDER_IMAGE_VERSION : "",
  ]
  cache-from = [
    LATEST_BUILDER_IMAGE,
    BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ,
  ]
  cache-to = [
    BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE
  ]
}
