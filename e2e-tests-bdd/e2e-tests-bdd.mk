LOCAL_DOCKER_COMPOSE_SERVICES += local-e2e-tests-bdd

APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE=$(CURDIR)/.containers/docker-buildx-cache/e2e-tests-bdd
$(APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE):;	@mkdir -p "$@"
ifneq ($(wildcard $(APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)/index.json),)
APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ = type=local,src=$(APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)
endif
APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE= type=local,dest=$(APPLICATION_TESTS_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)
APPLICATION_TESTS_IMAGE = ${IMAGES_BASE}:e2e-tests-bdd-latest

local-e2e-tests-bdd-docker-compose-up: local-e2e-tests-bdd-docker-compose-rm
	@VERSION=latest $(CURDIR)/e2e-tests-bdd/docker-compose-up.sh &
	@$(LOCAL_DOCKER_COMPOSE) logs e2e-tests-bdd
	@echo
	@echo Started Application tests display: 
	@echo Run "'make $(@:%-up=%-logs)'" to retrieve more logs

local-e2e-tests-bdd-docker-compose-run: 
	@$(LOCAL_DOCKER_COMPOSE) exec e2e-tests-bdd rf
.PHONY: e2e-tests-bdd-docker-compose-run

local-e2e-tests-bdd-check-env:
	@echo Checking application tests environment validity	
	$(LOCAL_DOCKER_COMPOSE) exec e2e-tests-bdd bash $(CURDIR)/application/tests/check_env.sh
.PHONY: local-e2e-tests-bdd-check-env
