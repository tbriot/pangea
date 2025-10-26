# ==============================================================================
# VARIABLES
# ==============================================================================

# Nom du programme principal et du port
PROGRAM_NAME = main
PORT = 8000

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

test-local: ## Exécute le test d'intégration via cURL.
	@echo "--------------------------------------------------------"
	@echo "🧪 Exécution du test d'intégration cURL..."
	@echo "--------------------------------------------------------"
	@echo ""
	@curl -iv -X POST http://127.0.0.1:8000/event \
		-H @test/01_push_event_headers \
		-d @test/01_push_event_payload.json

test-amzn-apigw: ## Exécute le test d'intégration via cURL.
	@echo "--------------------------------------------------------"
	@echo "🧪 Exécution du test d'intégration cURL..."
	@echo "--------------------------------------------------------"
	@echo ""
	@curl -iv -X POST ${AMZN_API_GATEWAY_URL} \
		-H @test/01_push_event_headers \
		-d @test/01_push_event_payload.json

test-amzn-ecs: ## Exécute le test d'intégration via cURL.
	@echo "--------------------------------------------------------"
	@echo "🧪 Exécution du test d'intégration cURL..."
	@echo "--------------------------------------------------------"
	@echo ""
	@curl -iv -X POST http://35.180.87.120:8000/event \
		-H @test/01_push_event_headers \
		--data-binary @test/01_push_event_payload.json

load-test-amzn-ecs: ## Exécute le test d'intégration via cURL.
	@echo "--------------------------------------------------------"
	@echo "🧪 Exécution du test d'intégration cURL..."
	@echo "--------------------------------------------------------"
	@echo ""
	@wrk "http://35.180.87.120:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 10 -d 30s

# ==============================================================================
# DOCKER TASKS
# ==============================================================================

docker-build: ## Construit l'image Docker.
	@echo "--------------------------------------------------------"
	@echo "🐳 Construction de l'image Docker $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	@echo "--------------------------------------------------------"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

ecr-login: 
	aws ecr get-login-password --region eu-west-3 | docker login --username AWS --password-stdin 590183718735.dkr.ecr.eu-west-3.amazonaws.com

docker-push:
	docker tag github-events-webhook:latest 590183718735.dkr.ecr.eu-west-3.amazonaws.com/github-events-webhook:latest
	docker push 590183718735.dkr.ecr.eu-west-3.amazonaws.com/github-events-webhook:latest

docker-run: ## Lance le conteneur Docker.
	@echo "--------------------------------------------------------"
	@echo "🚀 Démarrage du conteneur Docker $(DOCKER_IMAGE)..."
	@echo "--------------------------------------------------------"
	docker run -d --name $(DOCKER_IMAGE) -p $(PORT):$(PORT) \
		-e GITHUB_WEBHOOK_SECRET="$(GITHUB_WEBHOOK_SECRET)" \
		-e SQS_QUEUE_URL=https://sqs.eu-west-3.amazonaws.com/590183718735/github-webhook-events \
		-e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
		-e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
		-e AWS_SESSION_TOKEN=$(AWS_SESSION_TOKEN) \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

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
