package v1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

func shouldInjectSidecar(labels map[string]string) bool {
	if enable, _ := labels["inject-sidecar"]; enable == "enable" {
		return true
	}
	return false
}

func getSidecarContainer() corev1.Container {
	return corev1.Container{
		Name:  "demo-sidecar",
		Image: "nginx:1.23.3-alpine",
	}
}

func injectSidecar(pod *corev1.Pod) {
	sidecarContainer := getSidecarContainer()
	pod.Spec.Containers = append(pod.Spec.Containers, sidecarContainer)
	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string, 2)
	}
	pod.Annotations["sidecar_injected_at"] = time.Now().Format(`2006-01-02T15:04:05Z`)
	pod.Annotations["sidecar_injected_by"] = "webhook-kb-demo"
	// TODO 启动顺序控制
}
