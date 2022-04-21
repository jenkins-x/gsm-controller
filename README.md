# gsm-controller

[![Documentation](https://godoc.org/github.com/jenkins-x-labs/gsm-controller?status.svg)](https://pkg.go.dev/mod/github.com/jenkins-x-labs/gsm-controller)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x-labs/gsm-controller)](https://goreportcard.com/report/github.com/jenkins-x-labs/gsm-controller)
[![Releases](https://img.shields.io/github/release-pre/jenkins-x-labs/gsm-controller.svg)](https://github.com/jenkins-x-labs/gsm-controller/releases)
[![LICENSE](https://img.shields.io/github/license/jenkins-x-labs/gsm-controller.svg)](https://github.com/jenkins-x-labs/gsm-controller/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://slack.k8s.io/)

# Overview

gsm-controller is a Kubernetes controller that copies secrets from Google Secrets Manager into Kubernetes secrets.  The controller
watches Kubernetes secrets looking for an annotation, if the annotation is not found on the secret nothing more is done.

If the secret does have the annotation then the controller will query Google Secrets Manager, access the matching
secret and copy the value into the Kubernetes secret and save it in the cluster.

# Setup

_Note_ in this example we are creating secrets and running the Kubernetes cluster in the same Google Cloud Project, the same
approach will work if Secrets Manager is enabled in a different project to store your secrets, just set the env var `SECRETS_MANAGER_PROJECT_ID`
below to a different GCP project id.

Set some environment variables:
```bash
export NAMESPACE=foo
export CLUSTER_NAME=test-cluster-foo
export PROJECT_ID=my-cool-project
export SECRETS_MANAGER_PROJECT_ID=my-cool-project # change if you want you secrets stored in Secrets Manager from another GCP project
```

First enable Google Secrets Manager

```bash
gcloud services enable secretmanager.googleapis.com --project $SECRETS_MANAGER_PROJECT_ID
```

Create a secret
- Using a file:
```bash
gcloud beta secrets create foo --replication-policy automatic --project $SECRETS_MANAGER_PROJECT_ID --data-file=-=my_secrets.yaml
```
- or for a single key=value secret:
```bash
echo -n bar | gcloud beta secrets create foo --replication-policy automatic --project $SECRETS_MANAGER_PROJECT_ID --data-file=-
```


## Access

So that `gsm-controller` can access secrets in Google Secrets Manager so it can populate Kubernetes secrets in a namespace, it
requires a GCP service account with a role to access the secrets in a given GCP project.

### Setup
```bash
kubectl create serviceaccount gsm-sa -n $NAMESPACE
kubectl annotate sa gsm-sa iam.gke.io/gcp-service-account=$CLUSTER_NAME-sm@$SECRETS_MANAGER_PROJECT_ID.iam.gserviceaccount.com

gcloud iam service-accounts create $CLUSTER_NAME-sm --project $SECRETS_MANAGER_PROJECT_ID

gcloud iam service-accounts add-iam-policy-binding \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:$PROJECT_ID.svc.id.goog[$NAMESPACE/gsm-sa]" \
  $CLUSTER_NAME-sm@$SECRETS_MANAGER_PROJECT_ID.iam.gserviceaccount.com \
  --project $SECRETS_MANAGER_PROJECT_ID

gcloud projects add-iam-policy-binding $SECRETS_MANAGER_PROJECT_ID \
  --role roles/secretmanager.secretAccessor \
  --member "serviceAccount:$CLUSTER_NAME-sm@$SECRETS_MANAGER_PROJECT_ID.iam.gserviceaccount.com" \
  --project $SECRETS_MANAGER_PROJECT_ID
```

It can take a little while for permissions to propagate when using workload identity so it's a good idea to validate
auth is working before continuing to the next step.

run a temporary pod with our kubernetes service accounts

```bash
kubectl run --rm -it \
  --generator=run-pod/v1 \
  --image google/cloud-sdk:slim \
  --serviceaccount gsm-sa \
  --namespace $NAMESPACE \
  workload-identity-test
```
use gcloud to verify you can auth, it make take a few tries over a few minutes
```bash
gcloud auth list
```

# Install

add the helm repo or update it to get the latest charts
```bash
helm plugin install https://github.com/hayorov/helm-gcs
helm repo add jx3 https://jenkins-x-charts.github.io/repo
```
# or
```bash
helm repo update
```
install the helm chart, this includes a kubernetes controller that always runs and watches for new or updated secrets.  We also install a kubernetes CronJob that periodically triggers and checks for updated secret versions in Google Secret Manager.

```bash
helm install --set projectID=$SECRETS_MANAGER_PROJECT_ID gsm-controller jx3/gsm-controller
```

### Annotate secrets
Now that the controller is running we can create a Kubernetes secret and annotate it with the id we stored the secret
with above.

```bash
kubectl create secret generic my-secret
kubectl annotate secret my-secret jenkins-x.io/gsm-kubernetes-secret-key=credentials.json
kubectl annotate secret my-secret jenkins-x.io/gsm-secret-id=foo
```  
After a short wait you should be able to see the base64 encoded data in the secret
```bash
kubectl get secret my-secret -oyaml
```

If not check the logs of the controller
```bash
kubectl logs deployment/gsm-controller
```
### Run locally


```bash
gcloud iam service-accounts create $CLUSTER_NAME-sm --project $SECRETS_PROJECT_ID

gcloud iam service-accounts keys create ~/.secret/key.json \
  --iam-account $CLUSTER_NAME-sm@$PROJECT_ID.iam.gserviceaccount.com

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --role roles/secretmanager.secretAccessor \
  --member "serviceAccount:$CLUSTER_NAME-sm@$PROJECT_ID.iam.gserviceaccount.com"

```

Create a GCP secret in the project your secrets are stored, assign the accessor role, download the key.json and...
```bash
export GOOGLE_APPLICATION_CREDENTIALS=~/.secret/key.json
make build
./build/gsm-controller my-cool-project
```

# Realtime Updates

gsm-controller supports following Google Secret Manager's [event notifications](https://cloud.google.com/secret-manager/docs/event-notifications).

## Pubsub Access

```bash
gcloud beta services identity create \
    --service "secretmanager.googleapis.com" \
    --project $SECRETS_MANAGER_PROJECT_ID

# ^ This will print out a service account which can be constructed as follows ˅

PROJECT_NUMBER=gcloud projects describe $SECRETS_MANAGER_PROJECT_ID --format='value(projectNumber)'
export SM_SERVICE_ACCOUNT=service-${PROJECT_NUMBER}@gcp-sa-secretmanager.iam.gserviceaccount.com


gcloud pubsub topics create "projects/${SECRETS_MANAGER_PROJECT_ID}/topics/secrets.events"

gcloud pubsub topics add-iam-policy-binding secrets.events \
    --member "serviceAccount:${SM_SERVICE_ACCOUNT}" \
    --role "roles/pubsub.publisher" \
    --project $SECRETS_MANAGER_PROJECT_ID
```

Create a subscription for the Kubernetes cluster.
```bash
gcloud pubsub subscriptions create \
    "projects/${SECRETS_MANAGER_PROJECT_ID}/subscriptions/secrets.events.$CLUSTER_NAME.gsm-pubsub" \
    --topic "projects/${SECRETS_MANAGER_PROJECT_ID}/topics/secrets.events"
```

Allow reading the subscription
```bash
gcloud pubsub subscriptions add-iam-policy-binding secrets.events.$CLUSTER_NAME.gsm-pubsub \
    --member "serviceAccount:$CLUSTER_NAME-sm@$SECRETS_MANAGER_PROJECT_ID.iam.gserviceaccount.com" \
    --role "roles/pubsub.subscriber" \
    --project $SECRETS_MANAGER_PROJECT_ID
```

### Creating Secrets
When creating secrets, make sure to specify the topic
```bash
gcloud beta secrets create foo --replication-policy automatic \
    --project $SECRETS_MANAGER_PROJECT_ID --data-file=-=my_secrets.yaml \
    --topics "projects/${SECRETS_MANAGER_PROJECT_ID}/topics/secrets.events"
```


## Install with Pubsub

When following the [install instructions](#install), also set the `deployment.pubsub` values.
```bash
helm install \
  --set projectID=$SECRETS_MANAGER_PROJECT_ID \
  --set deployment.pubsub.enabled=true \
  --set deployment.pubsub.subscription=secrets.events.$CLUSTER_NAME.gsm-pubsub \
  --set cron.schedule="23 */2 * * *" \
  gsm-controller jx3/gsm-controller
```

The cron schedule of `"23 */2 * * *"` means the cronjob will run every 2 hours rather than every 5 minutes.

# Video

[![GSM Controller](http://img.youtube.com/vi/wLHgkhzeNe8/0.jpg)](http://www.youtube.com/watch?v=wLHgkhzeNe8 "Google Secrets Manager - Kubernetes Controller")

