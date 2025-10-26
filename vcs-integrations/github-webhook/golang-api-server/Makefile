# ==============================================================================
# VARIABLES
# ==============================================================================

# Nom du programme principal et du port
PROGRAM_NAME = main
PORT = 8000
API_ENDPOINT = http://localhost:$(PORT)/api/test
PAYLOAD = '{"message": "Hello Test"}'
HEADERS = -H 'Content-Type: application/json' -H 'Authorization: Bearer my-secret-token'

# Docker variables
DOCKER_IMAGE = github-events-webhook
DOCKER_TAG = latest

# Amazon API Gateway
AMZN_API_GATEWAY_URL = https://w2y8l5vnsg.execute-api.eu-west-3.amazonaws.com/test/event


# ==============================================================================
# TÂCHES PRINCIPALES
# ==============================================================================

.PHONY: all run test clean docker-build docker-run docker-stop docker-clean docker-scan

all: run ## Lance l'application et attend les requêtes.

run:
	@echo "--------------------------------------------------------"
	@echo "🚀 Démarrage du serveur sur $(API_ENDPOINT)"
	@echo "--------------------------------------------------------"
	go run $(PROGRAM_NAME).go

test: ## Exécute le test d'intégration via cURL.
	@echo "--------------------------------------------------------"
	@echo "🧪 Exécution du test d'intégration cURL..."
	@echo "--------------------------------------------------------"
	@echo ""
	@curl -iv -X POST http://127.0.0.1:8000 \
		-H @test/01_push_event_headers \
		-d @test/01_push_event_payload.json

test-aws: ## Exécute le test d'intégration via cURL.
	@echo "--------------------------------------------------------"
	@echo "🧪 Exécution du test d'intégration cURL..."
	@echo "--------------------------------------------------------"
	@echo ""
	@curl -iv -X POST ${AMZN_API_GATEWAY_URL} \
		-H @test/01_push_event_headers \
		-d @test/01_push_event_payload.json

# ==============================================================================
# DOCKER TASKS
# ==============================================================================

docker-build: ## Construit l'image Docker.
	@echo "--------------------------------------------------------"
	@echo "🐳 Construction de l'image Docker $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	@echo "--------------------------------------------------------"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Lance le conteneur Docker.
	@echo "--------------------------------------------------------"
	@echo "🚀 Démarrage du conteneur Docker $(DOCKER_IMAGE)..."
	@echo "--------------------------------------------------------"
	docker run -d --name $(DOCKER_IMAGE) -p $(PORT):$(PORT) $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop: ## Arrête le conteneur Docker.
	@echo "--------------------------------------------------------"
	@echo "🛑 Arrêt du conteneur Docker $(DOCKER_IMAGE)..."
	@echo "--------------------------------------------------------"
	docker stop $(DOCKER_IMAGE) || true
	docker rm $(DOCKER_IMAGE) || true

docker-clean: ## Nettoie les ressources Docker (image et conteneur).
	@echo "--------------------------------------------------------"
	@echo "🧹 Nettoyage des ressources Docker..."
	@echo "--------------------------------------------------------"
	docker stop $(DOCKER_IMAGE) || true
	docker rm $(DOCKER_IMAGE) || true
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true

docker-logs: ## Affiche les logs du conteneur Docker.
	docker logs -f $(DOCKER_IMAGE)

docker-scan: ## Scan l'image Docker avec Trivy pour détecter les vulnérabilités.
	@echo "--------------------------------------------------------"
	@echo "🔍 Scan de vulnérabilités de l'image Docker $(DOCKER_IMAGE):$(DOCKER_TAG) avec Trivy..."
	@echo "--------------------------------------------------------"
	trivy image $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-dev: docker-build docker-run ## Construit et lance l'image Docker.
