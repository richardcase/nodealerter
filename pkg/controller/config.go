package controller

import (
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Config represents the config for the conntroller
type Config struct {
	KubeConfig *rest.Config
	KubeClient kubernetes.Interface

	ResyncPeriod time.Duration

	NodesThreshold int
}
