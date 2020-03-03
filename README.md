# gsm-controller - experimental

[![Documentation](https://godoc.org/github.com/jenkins-x-labs/gsm-controller?status.svg)](https://pkg.go.dev/mod/github.com/jenkins-x-labs/gsm-controller)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x-labs/gsm-controller)](https://goreportcard.com/report/github.com/jenkins-x-labs/gsm-controller)
[![Releases](https://img.shields.io/github/release-pre/jenkins-x-labs/gsm-controller.svg)](https://github.com/jenkins-x-labs/gsm-controller/releases)
[![LICENSE](https://img.shields.io/github/license/jenkins-x-labs/gsm-controller.svg)](https://github.com/jenkins-x-labs/gsm-controller/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://slack.k8s.io/)

# Overview

gsm-controller is a Kubernetes controller that links Google Secrets Manager with Kubernetes secrets.  The controller
watches Kubernetes secrets looking for an annotation, if the annotation is not found on the secret nothing more is done.

If the secret does have the annotation then the controller will query Google Secrets Manager, access the matching
secret and copy teh value into the Kubernetes secret and save it in the cluster.

# Setup

First create a secret in Google Secrets Manager

Using a file:
```bash
gcloud beta secrets create foo --replication-policy automatic --project my-cool-project --data-file=-=my_secrets.yaml
```
or for a single key=value secret:
```bash
echo bar | gcloud beta secrets create foo --replication-policy automatic --project my-cool-project --data-file=-
```


## Access

So that `gsm-controller` can access secrets in Google Secrets Manager so it can populate Kubernetes secrets in a namespace, it
requires a GCP service account with a role to access the secrets in a given GCP project.

Set some environment variables:
```bash
export NAMESPACE=jx
export CLUSTER_NAME=test-cluster-foo
export PROJECT_ID=jx-development
```

### Setup
```bash
kubectl create serviceaccount gsm-sa
kubectl annotate sa gsm-sa jenkins-x.io/gsm-secret-id='foo'

gcloud iam service-accounts create $CLUSTER_NAME-sm

gcloud iam service-accounts add-iam-policy-binding \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:$PROJECT_ID.svc.id.goog[$NAMESPACE/gsm-sa]" \
  $CLUSTER_NAME-sm@$PROJECT_ID.iam.gserviceaccount.com

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --role roles/secretmanager.secretAccessor \
  --member "serviceAccount:$CLUSTER_NAME-sm@$PROJECT_ID.iam.gserviceaccount.com"
```

It can a little while for permissions to propogate when using workload identity so it's a good idea to validate auth is working before continuing to the next step.

run a temporary pod with one of our kubernetes service accounts

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

install the gsm controller chart
```bash
helm install gsm-controller \
  --set boot.namespace=$NAMESPACE \
  --set boot.projectID=$PROJECT_ID \
  .
```
### Run locally

Create a GCP secret in the project your secrets are stored, assign the accessor role, download the key.json and...
```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/credentials.json
make build
./build/gsm-controller my-cool-project
```




