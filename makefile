SERVICES ?= shared monitoring notification-service vanilla-server

S ?=
FILTERED_SERVICES := $(if $(S),$(filter $(S),$(SERVICES)),$(SERVICES))

NET_NAME ?= pet-project-first-network

.PHONY: help init copy-envs up stop restart network clean ps logs

help:
	@echo "Usage:"
	@echo "  make up [S=\"svc1 svc2\"]     - Start services (all or filtered)"
	@echo "  make stop [S=...]            - Stop services (all or filtered)"
	@echo "  make restart [S=...]         - Restart services"
	@echo "  make network                 - Ensure docker network exists ($(NET_NAME))"
	@echo "  make init                    - Create network, copy envs, start"
	@echo "  make copy-envs               - Copy .env.example -> .env where needed"
	@echo "  make logs [S=...]            - Proxy 'logs' into each service"
	@echo "  make ps   [S=...]            - Proxy 'ps' into each service"
	@echo "  make clean [S=...]           - Proxy 'down' into each service"
	@echo ""
	@echo "Tip: any other target will be proxied to sub-services too, e.g.:"
	@echo "  make migrate S=\"notification-service\""

init:
	@echo "🛠️  Initializing all services and network..."
	$(MAKE) network
	$(MAKE) copy-envs
	$(MAKE) up

copy-envs:
	@echo "📄 Copy .env.example -> .env (if missing)..."
	for service in $(FILTERED_SERVICES); do
		if [ -f $$service/.env.example ]; then
			if [ ! -f $$service/.env ]; then
				cp $$service/.env.example $$service/.env
				echo "✅ $$service: .env created"
			else
				echo "ℹ️  $$service: .env already exists, skipping"
			fi
		else
			echo "⚠️  $$service: .env.example not found"
		fi
	done

define CALL_IN_DIR
	if [ -f $(1)/Makefile ]; then \
		$(MAKE) -C $(1) $(2); \
	else \
		echo "⚠️  $(1): Makefile not found, skipping"; \
	fi
endef

up:
	@echo "🚀 Starting services: $(FILTERED_SERVICES)"
	@for service in $(FILTERED_SERVICES); do \
		echo "🟢 Starting $$service..."; \
		if [ -f $$service/Makefile ]; then \
			$(MAKE) -C $$service up || { echo "❌ Error on 'up' in $$service"; exit 1; }; \
			if [ "$$service" = "shared" ]; then \
				echo "⏳ Waiting 20s for $$service to be ready..."; \
				sleep 20; \
			fi; \
		else \
			echo "⚠️  $$service: Makefile not found, skipping"; \
		fi; \
	done

stop:
	@echo "🛑 Stopping services: $(FILTERED_SERVICES)"
	@for service in $(FILTERED_SERVICES); do \
		echo "🔴 Stopping $$service..."; \
		if [ -f $$service/Makefile ]; then \
			$(MAKE) -C $$service stop || echo "❗ Error on 'stop' in $$service"; \
		else \
			echo "⚠️  $$service: Makefile not found, skipping"; \
		fi; \
	done

restart:
	@$(MAKE) stop S="$(S)"
	@$(MAKE) up   S="$(S)"

network:
	@echo "🌐 Ensuring docker network '$(NET_NAME)' exists..."
	if ! docker network inspect $(NET_NAME) >/dev/null 2>&1; then
		docker network create --driver bridge --attachable $(NET_NAME)
		echo "✅ Network $(NET_NAME) created"
	else
		echo "ℹ️  Network $(NET_NAME) already exists"
	fi

logs:
	for service in $(FILTERED_SERVICES); do
		$(call CALL_IN_DIR,$$service,logs)
	done

ps:
	for service in $(FILTERED_SERVICES); do
		$(call CALL_IN_DIR,$$service,ps)
	done

clean:
	for service in $(FILTERED_SERVICES); do
		$(call CALL_IN_DIR,$$service,down)
	done

%:
	@set -e
	for service in $(FILTERED_SERVICES); do
		$(call CALL_IN_DIR,$$service,$@)
	done
