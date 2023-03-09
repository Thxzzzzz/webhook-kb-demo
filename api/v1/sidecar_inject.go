package v1

import (
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
		Env: []corev1.EnvVar{
			{
				Name:  "inject_from",
				Value: "webhook-kb-demo",
			},
		},
	}
}

func injectSidecar(pod *corev1.Pod) {
	sidecarContainer := getSidecarContainer()
	pod.Spec.Containers = append(pod.Spec.Containers, sidecarContainer)

	// TODO 启动顺序控制
}
