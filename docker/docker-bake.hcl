variable "VERSION" {
  description = "current git tag or commit version"
  default     = "dev"
}

variable "APPLICATION_DIR" {
  default = "application"
}

variable "NODE_VERSION" {
  description = "current Node version"
  default     = "22.14.0"
}

variable "NVM_VERSION" {
  description = "current NVM version"
  default     = "0.40.1"
}

variable "DOCKER_VERSION" {
  description = "current Docker version"
  default     = "28.0.2"
}

variable "DOCKER_COMPOSE_VERSION" {
  description = "current Docker Compose version"
  default     = "2.34.0"
}

variable "DOCKER_BUILDX_VERSION" {
  description = "current Docker Buildx version"
  default     = "0.22.0"
}

variable "TERRAFORM_VERSION" {
  description = "current Terraform version"
  default     = "1.11.2"
}

variable "KUBECTL_VERSION" {
  description = "current Terraform version"
  default     = "1.32"
}

variable "INFRA_ENV" {
  description = "production, staging or dev"
  default     = "dev"
}

variable "REPOSITORY" {
  default = "${INFRA_ENV}-docker-repository"
}

variable "PUBLIC_REPOSITORY" {
  default = "${INFRA_ENV}-docker-repository-public"
}

variable "GOOGLE_CLOUD_PROJECT_ID" {
  description = "Google Cloud Project ID"
  default     = "moov-dev-439608"
}

variable "PUBLIC_IMAGES_BASE" {
  default = "${PUBLIC_REPOSITORY}/${GOOGLE_CLOUD_PROJECT_ID}"
}

variable "PRIVATE_IMAGES_BASE" {
  default = "${REPOSITORY}/${GOOGLE_CLOUD_PROJECT_ID}"
}

variable "platforms" {
  default = [
    "linux/amd64",
    "linux/arm64"
  ]
}

group "services-load-local-arch" {
  targets = [
    "android-studio",
    "clarinet-devnet",
    "firebase-emulators",
    "github-action-runner",
    "vault-dev",
    "application-tests",
  ]
}

group "services-push-multi-arch" {
  targets = [
    "android-studio-ci",
    "clarinet-devnet-ci",
    "firebase-emulators-ci",
    "github-action-runner-ci",
    "vault-dev-ci",
    "application-tests-ci",
  ]
}
