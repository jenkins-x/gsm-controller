module github.com/jenkins-x-labs/gsm-controller

go 1.19

require (
	cloud.google.com/go v0.57.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/jenkins-x/jx-logging v0.0.3
	github.com/joho/godotenv v1.3.0
	github.com/magiconair/properties v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	google.golang.org/genproto v0.0.0-20200430143042-b979b6f78d84
	k8s.io/api v0.16.9
	k8s.io/apimachinery v0.16.10-beta.0
	k8s.io/client-go v0.16.9
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.1 // indirect
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jenkins-x/logrus-stackdriver-formatter v0.2.3 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/rickar/props v0.0.0-20170718221555-0b06aeb2f037 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.opencensus.io v0.22.3 // indirect
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79 // indirect
	golang.org/x/net v0.0.0-20200505041828-1ed23360d12c // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20200501145240-bc7a7d42d5c3 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/api v0.23.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	google.golang.org/protobuf v1.22.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20190801114015-581e00157fb1 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

replace github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
