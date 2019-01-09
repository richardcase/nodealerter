package controller

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) handlePodAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		logrus.Debugf("pod added: %s", key)
		c.queue.Add(key)
	}
}

func (c *Controller) handlePodDelete(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		logrus.Debugf("pod deleted: %s", key)
		c.queue.Add(key)
	}
}

func (c *Controller) handlePodUpdate(old, cur interface{}) {
	if old.(*corev1.Pod).GetResourceVersion() == cur.(*corev1.Pod).GetResourceVersion() {
		return
	}

	key, err := cache.MetaNamespaceKeyFunc(cur)
	if err == nil {
		logrus.Debugf("pod updated: %s", key)
		c.queue.Add(key)
	}
}
