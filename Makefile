BINARY_NAME=redis-counter

# service
redis-secret:
	kubectl create secret generic env-secrets --from-literal=REDIS_PASS=$(REDIS_PASS)

redis-disk:
	gcloud compute --project=$(GCP_PROJECT) disks create \
		redis-disk --zone=$(CLUSTER_ZONE) --type=pd-ssd --size=10GB

redis:
	kubectl apply -f deployments/redis-pd.yaml

# app
deps:
	go mod tidy

image:
	gcloud builds submit \
		--project $(GCP_PROJECT) \
		--tag gcr.io/$(GCP_PROJECT)/$(BINARY_NAME):latest

docker:
	docker build -t $(BINARY_NAME) .

deployment:
	kubectl apply -f deployments/app.yaml

cleanup:
	kubectl delete -f deployments/app.yaml