.PHONY: deploy
deploy: deploy-namespace deploy-elastic deploy-fluentd deploy-fluent-bit 

.PHONY: deploy-namespace
deploy-namespace:
	@kubectl apply -f ./namespace/*

.PHONY: deploy-fluent-bit
deploy-fluent-bit:
	@kubectl apply -f ./fluent/fluent_bit/*

.PHONY: deploy-fluentd
deploy-fluentd:
	@kubectl apply -f ./fluent/fluentd/*

.PHONY: deploy-elastic
deploy-elastic:
	@kubectl apply -f ./elastic_search

