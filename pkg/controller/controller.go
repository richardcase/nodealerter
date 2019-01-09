package controller

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Controller implements the node alerter controller
type Controller struct {
	Config

	kubeInformerFactory kubeinformers.SharedInformerFactory
	podsLister          corev1listers.PodLister

	queue workqueue.RateLimitingInterface

	hasSyncedFuncs []cache.InformerSynced
}

// New creates a new node alerter controller
func New(conf Config) *Controller {
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(conf.KubeClient, conf.ResyncPeriod)

	podInformer := kubeInformerFactory.Core().V1().Pods()

	controller := &Controller{
		Config:              conf,
		queue:               workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "nodealerter"),
		kubeInformerFactory: kubeInformerFactory,
		podsLister:          podInformer.Lister(),
	}

	controller.hasSyncedFuncs = []cache.InformerSynced{
		podInformer.Informer().HasSynced,
	}

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.handlePodAdd,
		DeleteFunc: controller.handlePodDelete,
		UpdateFunc: controller.handlePodUpdate,
	})

	return controller
}

func (c *Controller) Run(ctx context.Context) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	go c.kubeInformerFactory.Start(ctx.Done())

	logrus.Debug("waiting for caches to sync")
	if ok := cache.WaitForCacheSync(ctx.Done(), c.hasSyncedFuncs...); !ok {
		return errors.New("failed waiting for caches to sync")
	}

	logrus.Debug("starting worker")
	go wait.Until(c.worker, time.Second, ctx.Done())

	<-ctx.Done()
	return nil

}
