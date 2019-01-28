# kres (Knative Redis Event Source)

Redis-based queue event source for Knative

```shell
API --> Redis
        Redis --> Event Source
                  Event Source --> Service
                                   Service --> Firestore
```

> TODO: refactor last service to write back to Redis for status in API

## Setup

First, change these bellow variables to your own and export them in your environment

```shell
export GCP_PROJECT=project1
export GCP_ZONE=us-west1-a
export K8S_NAMESPACE=demo
export REDIS_PASS=secure-password-this-is-not
```

### Namespace and Config Maps

Next, ensure the namespace you defined above exists

```shell
kubectl create ns $K8S_NAMESPACE
```

Also, create a config map (`cm`) holding the GCP project ID

> Note, this is only required by the current Firestore client

```shell
kubectl create cm global-config --from-literal=GCP_PROJECT_ID=$GCP_PROJECT
```

### Install Receiving Service

For this demo we are going to install `events` service which expects Cloud Events v0.2 `POST`. This service will parse the submitted event and persists it to Firestore DB

```shell
kubectl apply -f config/service.yaml
```

> Once everything is validated, we can add `visibility: cluster-local` label to the service metadata to prevent it from being accessed from outside of the cluster.

```shell
labels:
    serving.knative.dev/visibility: cluster-local
```

### Redis

First, create a secret that will be shared between the service and the event source.

```shell
kubectl create secret generic redis-secrets --from-literal=REDIS_PASS=$REDIS_PASS -n $K8S_NAMESPACE
```

To enable redis queue persistence, we are also going to create a small SSD disks

> For best results, make sure the disk is provisioned in the same ZONE as your cluster

```shell
gcloud compute disks create redis-disk \
    --project=$GCP_PROJECT \
    --zone=$GCP_ZONE \
    --type=pd-ssd \
    --size=10GB
```

Now that we have the secret and disk configured, we can install redis service

```shell
kubectl apply -f config/redis.yaml
```

### Event Source

Now we are ready to create the redis queue events source. This source will trigger previously configured service whenever a new item is added to the redis queue

```shell
kubectl apply -f config/source.yaml
```

## Demo

Soon there will be a UI to trigger the redis queue event. For now connect to your redis container

```shell
kubectl exec -it redis-0 -n demo /bin/sh
```

Then inside of the container, access the redis CLI and authenticate

```shell
redis-cli
auth #<type in your password you defined in $REDIS_PASS >
```

You should see redis respond `OK`




## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.