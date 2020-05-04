module github.com/jenkins-x-labs/gsm-controller

go 1.12

require (
	cloud.google.com/go v0.53.0
	github.com/cloudflare/cfssl v0.0.0-20190409034051-768cd563887f
	github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96 // indirect
	github.com/evanphx/json-patch v4.2.0+incompatible // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/uuid v1.1.1 // indirect
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/jenkins-x/jx-logging v0.0.1
	github.com/joho/godotenv v1.3.0
	github.com/magiconair/properties v1.8.1
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0
	google.golang.org/genproto v0.0.0-20200212174721-66ed5ce911ce
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20191004102349-159aefb8556b
	k8s.io/apimachinery v0.0.0-20191004074956-c5d2f014d689
	k8s.io/client-go v11.0.1-0.20191004102930-01520b8320fc+incompatible
	k8s.io/code-generator v0.17.3 // indirect
	k8s.io/utils v0.0.0-20200229041039-0a110f9eb7ab // indirect
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

replace github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
