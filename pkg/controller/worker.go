package controller

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) worker() {
	logrus.Debug("starting worker")
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for c.processNextWorkItem() {
		}
	}()

	wg.Wait()

	logrus.Debug("exiting worker")
}

func (c *Controller) processNextWorkItem() bool {
	logrus.Debug("process next work item")

	key, quit := c.queue.Get()
	if quit {
		logrus.Debug("processing has quit")
		return false
	}
	defer c.queue.Done(key)

	err := c.handlePod(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(errors.Wrapf(err, "processing pod %q failed", key))
	c.queue.AddRateLimited(key)

	return true
}

func (c *Controller) handlePod(key string) error {
	//TODO: this is where we need to add our lofic
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return errors.Wrap(err, "splitting key")
	}
	pod, err := c.podsLister.Pods(namespace).Get(name)
	if err != nil {
		return errors.Wrapf(err, "getting pod %s", key)
	}

	logrus.Debugf("we have pod: %s created on %s", pod.Name, pod.CreationTimestamp.Time.String())

	return nil
}
