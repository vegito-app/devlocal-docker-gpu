variable "APPLICATION_TESTS_IMAGES_BASE" {
  default = "${PUBLIC_IMAGES_BASE}:application-tests"
}

variable "APPLICATION_TESTS_IMAGE_VERSION" {
  default = notequal("dev", VERSION) ? "${PUBLIC_IMAGES_BASE}:application-tests-${VERSION}" : ""
}

variable "LATEST_APPLICATION_TESTS_IMAGE" {
  default = "${APPLICATION_TESTS_IMAGES_BASE}-latest"
}

target "application-tests-ci" {
  args = {
    builder_image = LATEST_BUILDER_IMAGE
  }
  context    = "application-tests"
  dockerfile = "Dockerfile"
  tags = [
    notequal("", VERSION) ? APPLICATION_TESTS_IMAGE_VERSION : "",
    LATEST_APPLICATION_TESTS_IMAGE,
  ]
  cache-from = [
    # LATEST_BUILDER_IMAGE,
    LATEST_APPLICATION_TESTS_IMAGE,
  ]
  cache-to = [
    "type=inline",
  ]
}

variable "APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE" {
  description = "local write cache for tests image build"
}

variable "APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ" {
  description = "local read cache for tests image build (cannot be used before first write)"
}

target "application-tests" {
  context    = "application-tests"
  dockerfile = "Dockerfile"
  args = {
    builder_image = LATEST_BUILDER_IMAGE
  }
  tags = [
    notequal("", VERSION) ? APPLICATION_TESTS_IMAGE_VERSION : "",
    LATEST_APPLICATION_TESTS_IMAGE,
  ]
  cache-from = [
    APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ,
    BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ,
  ]
  cache-to = [
    APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE,
  ]
}
