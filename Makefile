COMPOSE_CMD = podman compose -f compose.local.yaml

local: local-clean local-start

local-clean:
	$(COMPOSE_CMD) down

local-start:
	$(COMPOSE_CMD) build --no-cache
	$(COMPOSE_CMD) up -d