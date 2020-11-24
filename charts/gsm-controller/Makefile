CHART_REPO := gs://jenkinsxio/charts
NAME := gsm-controller

build: clean
	rm -rf requirements.lock
	helm dependency build
	helm lint

install: clean build
	helm install . --name ${NAME}

upgrade: clean build
	helm upgrade ${NAME} .

delete:
	helm delete --purge ${NAME}

clean:
	rm -rf charts
	rm -rf ${NAME}*.tgz

release: clean
	sed -i -e "s/version:.*/version: $(VERSION)/" Chart.yaml
	sed -i -e "s|repository:.*|repository: gcr.io\/jenkinsxio-labs\/$(NAME)|" values.yaml
	sed -i -e "s/tag:.*/tag: $(VERSION)/" values.yaml

	helm dependency build
	helm lint
	helm package .
	helm repo add jx-labs $(CHART_REPO)
	helm gcs push ${NAME}*.tgz jx-labs --public
	rm -rf ${NAME}*.tgz%
