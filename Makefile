
.PHONY: event

# REDIS

redis-secret:
	kubectl create secret generic redis-secrets \
		--from-literal=REDIS_PASS=$(REDIS_PASS)

redis-disk:
	gcloud compute --project=$(GCP_PROJECT) disks create \
		redis-disk --zone=$(CLUSTER_ZONE) --type=pd-ssd --size=10GB

redis:
	kubectl apply -f config/redis.yaml

forward:
	kubectl port-forward pods/redis-0 6379:6379 -n demo

# DEV

deps:
	go mod tidy

run:
	go run cmd/service/*.go --sink=https://events.demo.knative.tech/

image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/kres:latest

docker:
	docker build -t kres .

source:
	kubectl apply -f config/source.yaml

cleanup:
	kubectl delete -f config/source.yaml

event:
	curl -H "Content-Type: application/json" -X POST \
			"https://kapi.demo.knative.tech/v1/stock/oooo" | jq "."

status:
	curl -H "Content-Type: application/json" -X GET \
			"https://kapi.demo.knative.tech/v1/status/id-0a580a14019e"
