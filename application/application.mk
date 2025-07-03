-include application/go.mk
-include application/nodejs.mk
-include application/frontend/frontend.mk
-include application/backend/backend.mk
-include application/images/images.mk
# -include application/mobile/mobile.mk
-include application/tests/tests.mk
# -include application/run/run.mk

APPLICATION_CI_IMAGES_BUILD := \
	local-application-backend-image-ci \
	local-application-images-moderator-image-ci \
	local-application-images-cleaner-image-ci \
	local-application-tests-image-ci

application-images-ci: $(APPLICATION_CI_IMAGES_BUILD)
.PHONY: application-images-ci