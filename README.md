# kres (Knative Redis Event Source)

Redis-based queue event source for Knative

## Setup

Change the bellow variables to what you prefer and export them in your environment

```shell
export GCP_PROJECT=s9-demo
export GCP_ZONE=us-west1-a
export K8S_NAMESPACE=demo
export REDIS_PASS=longpasswordthisisinotsecure
```

### Namespace and Config Maps

Ensure your namespace exists

```shell
kubectl create ns $K8S_NAMESPACE
```

Also, create a config map holding the GCP project ID

> Note, this is still required with the current Firestore client

```shell
kubectl create cm global-config --from-literal=GCP_PROJECT_ID=$GCP_PROJECT
```

### Install Receiving Service

For this demo we are going to install `events` service which expects Cloud Events v0.2 in a form of a `POST` and persists them to a Firestore DB

```shell
kubectl apply -f config/service.yaml
```

### Redis

Create a secret

```shell
kubectl create secret generic redis-secrets --from-literal=REDIS_PASS=$REDIS_PASS -n $K8S_NAMESPACE
```

Create persistent disks

> For best results, make sure the disk is provisioned in the same ZONE as your cluster

```shell
gcloud compute --project=$GCP_PROJECT disks create redis-disk --zone=$GCP_ZONE --type=pd-ssd --size=10GB
```

Install redis service

```shell
kubectl apply -f config/redis.yaml
```

### Event Source

Now we are ready to create the redis queue events source.

```shell
kubectl apply -f config/source.yaml
```

## Demo

Connect to your redis container

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