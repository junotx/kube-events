module github.com/kubesphere/kube-events

go 1.13

require (
	github.com/antlr/antlr4 v0.0.0-20200225173536-225249fdaef5 // indirect
	github.com/go-logr/logr v0.1.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/kubesphere/alertmanager-kit v0.0.0-20201019060038-52e1f8a13968
	github.com/kubesphere/event-rule-engine v0.0.0-20200808103159-763922656585
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/panjf2000/ants/v2 v2.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.2.1 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.5.0
	sigs.k8s.io/controller-tools v0.4.1 // indirect
	sigs.k8s.io/kustomize/kustomize/v3 v3.10.0 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace k8s.io/client-go => k8s.io/client-go v0.17.2
