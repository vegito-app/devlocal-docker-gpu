CLARINET_DEVNET_IMAGE_DOCKER_BUILDX_LOCAL_CACHE=$(CURDIR)/.containers/docker-buildx-cache/clarinet-devnet
$(CLARINET_DEVNET_IMAGE_DOCKER_BUILDX_LOCAL_CACHE):;	@mkdir -p "$@"
ifneq ($(wildcard $(CLARINET_DEVNET_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)/index.json),)
CLARINET_DEVNET_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ = type=local,src=$(CLARINET_DEVNET_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)
endif
CLARINET_DEVNET_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE= type=local,dest=$(CLARINET_DEVNET_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)

CONTRACTS := \
  my-first-contract \
  counter \
  order
BITCOIND_DOCKER_NETWORK_NAME = my-first-contract.devnet
BITCOIND_DOCKER_CONTAINER_NAME = bitcoin-node.$(BITCOIND_DOCKER_NETWORK_NAME)

# Variables
CONTRACT_DIR   = $(CURDIR)/application/contracts
DEPLOYMENT_DIR = $(CONTRACT_DIR)/deployment
CLARINET_DIR   = $(CURDIR)/clarinet-devnet
CLARINET       = clarinet 

local-clarinet-check: ## Vérifie la validité des contrats
	@echo "🔍 Vérification des contrats Clarity..."
	@for contract in $(CONTRACT_DIR)/*.clar; do \
		$(CLARINET) check $$contract || exit 1; \
	done
	@echo "✅ Tous les contrats sont valides."
.PHONY: local-clarinet-check 

local-clarinet-deploy-devnet: ## Déploie les contrats sur le Devnet
	@echo "🚀 Déploiement sur le Devnet..."
	@clarity-cli deploy --config $(DEPLOYMENT_DIR)/devnet.yaml
	@echo "✅ Déploiement terminé sur le Devnet."
.PHONY: local-clarinet-deploy-devnet 

local-clarinet-deploy-staging: ## Déploie les contrats sur le Staging
	@echo "🚀 Déploiement sur Staging..."
	@clarity-cli deploy --config $(DEPLOYMENT_DIR)/staging.yaml
	@echo "✅ Déploiement terminé sur Staging."
.PHONY: local-clarinet-deploy-staging 

local-clarinet-clean: ## Supprime les artefacts générés
	@echo "🧹 Nettoyage..."
	@rm -f $(CONTRACT_DIR)/*.wasm
	@echo "✅ Nettoyage terminé."
.PHONY: local-clarinet-clean 

local-clarinet-devnet-start:
	@-docker rm -f `docker ps -aq --filter name=devnet` 2>/dev/null
	@-docker network rm -f `docker network ls -q --filter name=devnet` 2>/dev/null
	@cd $(CLARINET_DIR) && clarinet devnet start --no-dashboard
.PHONY: local-clarinet-devnet-start

local-clarinet-devnet-docker-compose-up: local-clarinet-devnet-docker-compose-rm
	@$(CLARINET_DIR)/docker-compose-up.sh &
	@$(LOCAL_DOCKER_COMPOSE) logs clarinet-devnet
	@echo
	@echo Started Clarinet Devnet: 
	@echo Run "'make $(@:%-up=%-logs)'" to retrieve more logs
.PHONY: local-clarinet-devnet-docker-compose-up