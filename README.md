# knative-redis-event-source
Redis-based queue event source for Knative




## Install

### Redis


Create a secret

```shell
kubectl create secret generic redis-secrets --from-literal=REDIS_PASS=YOUR_SECRET_HERE
```

Create persistent disks

> For best results, make sure the disk is provisioned in the same ZONE as your cluster

```shell
gcloud compute --project=GCP_PROJECT disks create redis-disk --zone=GCP_ZONE --type=pd-ssd --size=10GB
```

Install redis

```shell
kubectl apply -f deployments/redis.yaml
```


To test, connect to newly created service

```shell
kubectl exec -it redis-0 /bin/sh
# redis-cli
# AUTH <yourpassword>
```

### Knative app

To deploy the `counter` app on Knative we are going to:

- [knative-redis-event-source](#knative-redis-event-source)
  - [Install](#install)
    - [Redis](#redis)
    - [Knative app](#knative-app)
      - [Build the image](#build-the-image)
      - [Configure Knative](#configure-knative)
      - [Deploy Service](#deploy-service)
  - [Test](#test)
  - [Disclaimer](#disclaimer)

#### Build the image

Quickest way to build your service image is through [GCP Build](https://cloud.google.com/cloud-build/). Just submit the build request from within the `gauther` directory:

```shell
gcloud builds submit \
    --project ${GCP_PROJECT} \
	--tag gcr.io/${GCP_PROJECT}/redis-counter:latest
```

The build service is pretty verbose in output but eventually you should see something like this

```shell
ID           CREATE_TIME          DURATION  SOURCE                                   IMAGES                      STATUS
6905dd3a...  2018-12-23T03:48...  1M43S     gs://PROJECT_cloudbuild/source/15...tgz  gcr.io/PROJECT/redis-counter SUCCESS
```

Copy the image URI from `IMAGE` column (e.g. `gcr.io/GCP_PROJECT/redis-counter`).

#### Configure Knative

Before we can deploy that service to Knative, we just need to update the `deployments/app.yaml` file to update the image URL:


```yaml
    spec:
        container:
            image: gcr.io/GCP_PROJECT/redis-counter:latest
```

#### Deploy Service

Once done updating our service manifest (`deployments/app.yaml`) you are ready to deploy it.

```shell
kubectl apply -f deployments/app.yaml
```

The response should be

```shell
service.serving.knative.dev "counter" configured
```

To check if the service was deployed successfully you can check the status using `kubectl get pods` command. The response should look something like this (e.g. Ready `3/3` and Status `Running`).

```shell
NAME                                          READY     STATUS    RESTARTS   AGE
counter-00001-deployment-5645f48b4d-mb24j      3/3       Running   0          4h
```

You should be able to test the app now in browser

## Test

First access the root of the app to make sure the application is deployed (for example, if your cluster is configured with the domain `knative.tech` the URL would look like this)

https://counter.default.knative.tech

If you get `OK` reponse you are ready to test the redis connection but accessing the `counter`

https://counter.default.knative.tech/count

The response should be

```json
{
    "counter": 18,
    "error": ""
}
```

Where the `/counter` field increments by 1 on each refresh

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.