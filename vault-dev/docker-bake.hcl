variable "VAULT_DEV_IMAGE_VERSION" {
  default = notequal("dev", VERSION) ? "${PUBLIC_IMAGES_BASE}:vault-dev-${VERSION}" : ""
}

variable "LATEST_VAULT_DEV_IMAGE" {
  default = "${PUBLIC_IMAGES_BASE}:vault-dev-latest"
}

target "vault-dev-ci" {
  context    = "vault-dev"
  dockerfile = "Dockerfile"
  tags = [
    LATEST_VAULT_DEV_IMAGE,
    notequal("", VERSION) ? VAULT_DEV_IMAGE_VERSION : "",
  ]
  cache-from = [
    LATEST_VAULT_DEV_IMAGE
  ]
  cache-to  = ["type=inline"]
  platforms = platforms
}

variable "VAULT_DEV_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE" {
  description = "local write cache for builder image build"
}

variable "VAULT_DEV_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ" {
  description = "local read cache for builder image build (cannot be used before first write)"
}

target "vault-dev" {
  context    = "vault-dev"
  dockerfile = "Dockerfile"
  tags = [
    LATEST_VAULT_DEV_IMAGE,
    notequal("", VERSION) ? VAULT_DEV_IMAGE_VERSION : "",
  ]
  cache-from = [
    LATEST_VAULT_DEV_IMAGE,
    VAULT_DEV_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ,
  ]
  cache-to = [
    VAULT_DEV_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE
  ]
}
