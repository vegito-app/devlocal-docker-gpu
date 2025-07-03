variable "FIREBASE_EMULATORS_IMAGE_TAG" {
  default = notequal("", VERSION) ? "${PUBLIC_IMAGES_BASE}:firebase-emulators-${VERSION}" : ""
}

variable "LATEST_FIREBASE_EMULATORS_IMAGE" {
  default = "${PUBLIC_IMAGES_BASE}:firebase-emulators-latest"
}

variable "FIREBASE_EMULATORS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE" {
  description = "local write cache for firebase-emulators image build"
}

variable "FIREBASE_EMULATORS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ" {
  description = "local read cache for firebase-emulators image build (cannot be used before first write)"
}

target "firebase-emulators-ci" {
  args = {
    builder_image = LATEST_BUILDER_IMAGE
  }
  context = "firebase-emulators"
  tags = [
    LATEST_FIREBASE_EMULATORS_IMAGE,
    FIREBASE_EMULATORS_IMAGE_TAG,
  ]
  cache-from = [
    LATEST_BUILDER_IMAGE,
    LATEST_FIREBASE_EMULATORS_IMAGE
  ]
  cache-to  = ["type=inline"]
  platforms = platforms
}

target "firebase-emulators" {
  args = {
    builder_image = LATEST_BUILDER_IMAGE
  }
  context = "firebase-emulators"
  tags = [
    LATEST_FIREBASE_EMULATORS_IMAGE,
    FIREBASE_EMULATORS_IMAGE_TAG,
  ]
  cache-from = [
    FIREBASE_EMULATORS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ,
    BUILDER_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ,
  ]
  cache-to = [
    FIREBASE_EMULATORS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE
  ]
}
