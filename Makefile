
.PHONY: client

# REDIS

redis-secret:
	kubectl create secret generic redis-secrets \
		--from-literal=REDIS_PASS=$(REDIS_PASS)

redis-disk:
	gcloud compute --project=$(GCP_PROJECT) disks create \
		redis-disk --zone=$(CLUSTER_ZONE) --type=pd-ssd --size=10GB

redis:
	kubectl apply -f config/redis.yaml

# DEV

deps:
	go mod tidy

image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/kres:latest

source:
	kubectl apply -f config/source.yaml

cleanup:
	kubectl delete -f config/source.yaml

client:
	go build ./cmd/client/