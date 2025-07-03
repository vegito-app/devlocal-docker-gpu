variable "GITHUB_RUNNER_IMAGE_VERSION" {
  default = notequal("dev", VERSION) ? "${PUBLIC_IMAGES_BASE}:github-action-runner-${VERSION}" : ""
}

variable "LATEST_GITHUB_RUNNER_IMAGE" {
  default = "${PUBLIC_IMAGES_BASE}:github-action-runner-latest"
}

variable "GITHUB_ACTION_RUNNER_VERSION" {
  description = "current Github Actions Runner version"
  default     = "2.323.0"
}

group "service" {
  targets = ["github-action-runner"]
}

group "local-service" {
  targets = ["github-action-runner-local"]
}

target "github-action-runner-ci" {
  args = {
    docker_version         = DOCKER_VERSION
    docker_compose_version = DOCKER_COMPOSE_VERSION
    docker_buildx_version  = DOCKER_BUILDX_VERSION
    terraform_version      = TERRAFORM_VERSION
    kubectl_version        = KUBECTL_VERSION
    node_version           = NODE_VERSION
    nvm_version            = NVM_VERSION
    github_runner_version  = GITHUB_ACTION_RUNNER_VERSION
  }
  depends_on = [builder]
  context    = "github"
  tags = [
    LATEST_GITHUB_RUNNER_IMAGE,
    notequal("", VERSION) ? GITHUB_RUNNER_IMAGE_VERSION : "",
  ]
  cache-from = [
    LATEST_BUILDER_IMAGE,
    LATEST_GITHUB_RUNNER_IMAGE
  ]
  cache-to  = ["type=inline"]
  platforms = platforms
}

variable "GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE" {
  description = "local write cache for github-actions-runner image build"
}

variable "GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ" {
  description = "local read cache for github-actions-runner image build (cannot be used before first write)"
}

target "github-action-runner" {
  args = {
    docker_version         = DOCKER_VERSION
    docker_compose_version = DOCKER_COMPOSE_VERSION
    docker_buildx_version  = DOCKER_BUILDX_VERSION
    terraform_version      = TERRAFORM_VERSION
    kubectl_version        = KUBECTL_VERSION
    node_version           = NODE_VERSION
    nvm_version            = NVM_VERSION
    github_runner_version  = GITHUB_ACTION_RUNNER_VERSION
  }
  depends_on = [builder-local]
  context    = "github"
  tags = [
    LATEST_GITHUB_RUNNER_IMAGE,
    notequal("", VERSION) ? GITHUB_RUNNER_IMAGE_VERSION : "",
  ]
  cache-from = [
    GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ,
    BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ,
  ]
  cache-to = [
    GITHUB_ACTIONS_RUNNER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE,
  ]
}
