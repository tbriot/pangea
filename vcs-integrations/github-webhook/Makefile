# ==============================================================================
# VARIABLES
# ==============================================================================

# Nom du programme principal et du port
PROGRAM_NAME = main
PORT = 8000
API_ENDPOINT = http://localhost:$(PORT)/api/test
PAYLOAD = '{"message": "Hello Test"}'
HEADERS = -H 'Content-Type: application/json' -H 'Authorization: Bearer my-secret-token'


# ==============================================================================
# TÂCHES PRINCIPALES
# ==============================================================================

.PHONY: all run test clean

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
