
.PHONY: help \
	setup init \
	generate-ca generate-server-cert generate-certs verify-certs \
	generate-redis-conf generate-acl \
	up down restart ps logs \
	clean reset prune

PROJECT_NAME := redis-security
COMPOSE_FILE := docker/redis.docker-compose.yml

DOCKER_DIR := docker
DATA_DIR := $(DOCKER_DIR)/redis-data
CERTS_DIR := $(DOCKER_DIR)/certs

REDIS_CONF := $(DOCKER_DIR)/redis.conf
ACL_FILE := $(DOCKER_DIR)/users.acl

CA_KEY := $(CERTS_DIR)/ca.key
CA_CERT := $(CERTS_DIR)/ca.crt

SERVER_KEY := $(CERTS_DIR)/server.key
SERVER_CSR := $(CERTS_DIR)/server.csr
SERVER_CERT := $(CERTS_DIR)/server.crt
SERVER_EXT := $(CERTS_DIR)/server.ext

# ======================== Help ========================

help:
	@echo ""
	@echo "Redis ACL + TLS Demo"
	@echo ""
	@echo "Setup"
	@echo "  make setup"
	@echo ""
	@echo "Redis"
	@echo "  make up"
	@echo "  make down"
	@echo "  make restart"
	@echo "  make ps"
	@echo "  make logs"
	@echo ""
	@echo "TLS"
	@echo "  make verify-certs"
	@echo ""
	@echo "Cleanup"
	@echo "  make clean"
	@echo "  make reset"
	@echo "  make prune"
	@echo ""

# ======================== Setup ========================

setup: init generate-certs generate-redis-conf generate-acl
	@echo "✅ Environment ready"

init:
	@echo "📁 Preparing directories..."
	@mkdir -p $(DATA_DIR) $(CERTS_DIR)
	@chmod 755 $(DATA_DIR) $(CERTS_DIR)
	@echo "✅ Directories ready"

# ======================== TLS ========================

generate-certs: generate-ca generate-server-cert
	@echo "✅ TLS certificates generated"

generate-ca:
	@echo "🔐 Generating Certificate Authority..."
	@openssl genrsa -out $(CA_KEY) 4096
	@openssl req -x509 -new -nodes \
		-key $(CA_KEY) \
		-sha256 \
		-days 365 \
		-out $(CA_CERT) \
		-subj "/CN=redis-ca"
	@chmod 600 $(CA_KEY)
	@chmod 644 $(CA_CERT)
	@echo "✅ CA generated"

generate-server-cert:
	@echo "🔐 Generating Redis server certificate..."
	@openssl genrsa -out $(SERVER_KEY) 4096
	@openssl req -new \
		-key $(SERVER_KEY) \
		-out $(SERVER_CSR) \
		-subj "/CN=localhost"

	@printf "subjectAltName=DNS:localhost,IP:127.0.0.1\n" > $(SERVER_EXT)

	@openssl x509 -req \
		-in $(SERVER_CSR) \
		-CA $(CA_CERT) \
		-CAkey $(CA_KEY) \
		-CAcreateserial \
		-out $(SERVER_CERT) \
		-days 365 \
		-sha256 \
		-extfile $(SERVER_EXT)

	@chmod 600 $(SERVER_KEY)
	@chmod 644 $(SERVER_CERT)

	@echo "✅ Server certificate generated"

verify-certs:
	@echo "🔍 Verifying certificates..."
	@openssl verify -CAfile $(CA_CERT) $(SERVER_CERT)
	@echo "✅ Certificates verified"

# ======================== Redis Configuration ========================

generate-redis-conf:
	@echo "📝 Generating redis.conf..."
	@{ \
		echo "appendonly yes"; \
		echo ""; \
		echo "port 0"; \
		echo "tls-port 6379"; \
		echo ""; \
		echo "tls-cert-file /certs/server.crt"; \
		echo "tls-key-file /certs/server.key"; \
		echo "tls-ca-cert-file /certs/ca.crt"; \
		echo "tls-auth-clients no"; \
		echo ""; \
		echo "aclfile /usr/local/etc/redis/users.acl"; \
	} > $(REDIS_CONF)
	@echo "✅ redis.conf generated"


generate-acl:
	@echo "📝 Generating users.acl..."
	@{ \
		echo "user default off"; \
		echo ""; \
		echo "user cache-service on >cache123 ~cache:* &cache:* +@read +@write +publish +subscribe"; \
	} > $(ACL_FILE)
	@echo "✅ users.acl generated"



# ======================== Redis ========================

up:
	@echo "🐳 Starting Redis..."
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) up -d --force-recreate 
	@echo "✅ Redis started"

down:
	@echo "🛑 Stopping Redis..."
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) down
	@echo "✅ Redis stopped"

restart: down up

ps:
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) ps

logs:
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) logs -f redis

# ======================== Cleanup ========================

clean:
	@echo "🧹 Removing Redis data..."
	@rm -rf $(DATA_DIR)
	@echo "✅ Redis data removed"

reset: down
	@echo "🧹 Resetting environment..."
	@rm -rf $(DATA_DIR)
	@rm -rf $(CERTS_DIR)
	@rm -f $(REDIS_CONF)
	@rm -f $(ACL_FILE)
	@echo "✅ Environment reset"

prune:
	@echo "⚠️ Running Docker system prune..."
	@docker system prune -f
	@echo "✅ Docker cleaned"

