module github.com/jenkins-x-labs/gsm-controller

go 1.12

require (
	cloud.google.com/go v0.56.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.1 // indirect
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/jenkins-x/jx-logging v0.0.3
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/spf13/cobra v1.0.0
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79 // indirect
	golang.org/x/net v0.0.0-20200505041828-1ed23360d12c // indirect
	golang.org/x/sys v0.0.0-20200501145240-bc7a7d42d5c3 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/api v0.23.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200430143042-b979b6f78d84
	google.golang.org/grpc v1.29.1 // indirect
	k8s.io/api v0.16.9
	k8s.io/apimachinery v0.16.10-beta.0
	k8s.io/client-go v0.16.9
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

replace github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
