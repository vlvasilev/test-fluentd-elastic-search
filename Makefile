REGISTRY                           := hisshadow85
SERVER_ONE_IMAGE_REPOSITORY        := $(REGISTRY)/server-one
FLOOD_AND_ANALYZE_IMAGE_REPOSITORY := $(REGISTRY)/flood-and-analyze
IMAGE_TAG                          := $(shell cat VERSION)

SHOOTS := $(if $(SHOOTS),$(SHOOTS),1)
PODS := $(if $(PODS),$(PODS),10)
MESSAGES := $(if $(MESSAGES),$(MESSAGES),10)
LOGGING_TIME := $(if $(LOGGING_TIME),$(LOGGING_TIME),30000)
TTWAL := $(if $(TTWAL),$(TTWAL),60)
MASTER := $(if $(MASTER),$(MASTER),"localhost:8000")


.PHONY: build-server-one
build-server-one:
	@CGO_ENABLED=0 & GOOS=linux & go build  -ldflags "-linkmode external -extldflags -static" -a -o ./bin/server_one/server_one ./cmd/server/server_one/server.go 2>/dev/null

.PHONY: build-flood-and-analyze
build-flood-and-analyze:
	@CGO_ENABLED=0 & GOOS=linux & go build  -ldflags "-linkmode external -extldflags -static" -a -o ./bin/flood_and_analyze/flood_and_analyze ./cmd/flood_and_analyze/flood_and_analyze.go 2>/dev/null

.PHONY: build
build: build-server-one build-flood-and-analyze

.PHONY: docker-build-server-one
docker-build-server-one:
	@if [ ! -f ./bin/server_one/server_one ]; then echo "No binary found. Please run 'make build-server-one' or 'make build'"; false; fi
	@docker build -t $(SERVER_ONE_IMAGE_REPOSITORY):$(IMAGE_TAG) -t $(SERVER_ONE_IMAGE_REPOSITORY):latest -f build/server_one/Dockerfile --rm .

.PHONY: docker-build-flood-and-analyze
docker-build-flood-and-analyze:
	@if [ ! -f ./bin/flood_and_analyze/flood_and_analyze ]; then echo "No binary found. Please run 'make build-flood-and-analyze' or 'make build'"; false; fi
	@docker build -t $(FLOOD_AND_ANALYZE_IMAGE_REPOSITORY):$(IMAGE_TAG) -t $(FLOOD_AND_ANALYZE_IMAGE_REPOSITORY):latest -f build/flood_and_analyze/Dockerfile --rm .

.PHONY: docker-build
docker-build: docker-build-server-one docker-build-flood-and-analyze

.PHONY: docker-push-server-one
docker-push-server-one:
	@if ! docker images $(SERVER_ONE_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(IMAGE_TAG); then echo "$(SERVER_ONE_IMAGE_REPOSITORY) version $(IMAGE_TAG) is not yet built. Please run 'make docker-build-server-one' or 'make docker-build'"; false; fi
	@docker push $(SERVER_ONE_IMAGE_REPOSITORY):$(IMAGE_TAG)
	@docker push $(SERVER_ONE_IMAGE_REPOSITORY):latest

.PHONY: docker-push-flood-and-analyze
docker-push-flood-and-analyze:
	@if ! docker images $(FLOOD_AND_ANALYZE_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(IMAGE_TAG); then echo "$(FLOOD_AND_ANALYZE_IMAGE_REPOSITORY) version $(IMAGE_TAG) is not yet built. Please run 'make docker-build-server-one' or 'make docker-build'"; false; fi
	@docker push $(FLOOD_AND_ANALYZE_IMAGE_REPOSITORY):$(IMAGE_TAG)
	@docker push $(FLOOD_AND_ANALYZE_IMAGE_REPOSITORY):latest

.PHONY: docker-push
docker-push: docker-push-server-one docker-push-flood-and-analyze

.PHONY: deploy-server-one
deploy-server-one:
	@kubectl apply -f ./docs/server/deploying/yaml/namespace.yaml
	@kubectl apply -f ./docs/server/deploying/yaml/serviceaccount.yaml
	@kubectl apply -f ./docs/server/deploying/yaml/clusterrolebinding.yaml
	@kubectl apply -f ./docs/server/deploying/yaml/deployment.yaml
	@kubectl apply -f ./docs/server/deploying/yaml/loadbalancer.yaml

.PHONY: deploy-center-stack
deploy-center-stack:
	@bash -x ./appliance/central_stack/deploy.sh

.PHONE: deploy-shoot-stack
deploy-shoot-stack:
	@./appliance/shoot_stack/deploy.sh $(SHOOTS)

.PHONY: deploy-appliance
deploy-appliance: deploy-center-stack deploy-shoot-stack deploy-server-one

.PHONY: run-test
run-test:
	@./appliance/tests/test_appliance.sh $(PODS) $(MESSAGES) $(LOGGING_TIME) $(TTWAL) $(MASTER)

.PHONY: clean
clean:
	@rm -rf bin/
