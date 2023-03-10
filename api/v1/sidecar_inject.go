package v1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	labelKeyInjectSidecar            = "inject-sidecar"
	labelKeyHoldAppUntilSidecarReady = "hold-app-until-sidecar-ready"
)

func shouldInjectSidecar(labels map[string]string) bool {
	if enable, _ := labels[labelKeyInjectSidecar]; enable == "enable" {
		return true
	}
	return false
}

func shouldHoldAppUntilSidecarReady(labels map[string]string) bool {
	if enable, _ := labels[labelKeyHoldAppUntilSidecarReady]; enable == "enable" {
		return true
	}
	return false
}

func getSidecarContainer(holdAppUntilProxyReady bool) (sidecar corev1.Container, injectFront bool) {
	sidecar = corev1.Container{
		Name:  "demo-sidecar",
		Image: "379809513/demo-sidecar:latest",
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{Path: "/healthz", Port: intstr.FromInt(9000)},
			},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{Path: "/readyz", Port: intstr.FromInt(9000)},
			},
		},
		ImagePullPolicy: corev1.PullAlways,
	}

	if !holdAppUntilProxyReady {
		return
	}
	injectFront = true
	sidecar.Lifecycle = &corev1.Lifecycle{
		PostStart: &corev1.LifecycleHandler{
			Exec: &corev1.ExecAction{
				// TODO wait until ready
				Command: []string{
					"/bin/sh", "-c",
					`while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:9000/healthz)" != '200' ]]; do echo Waiting for Sidecar;sleep 1; done; echo Sidecar available;`,
				},
			},
		},
	}
	return
}

func injectSidecar(pod *corev1.Pod) {
	sidecarContainer, injectFront := getSidecarContainer(shouldHoldAppUntilSidecarReady(pod.Labels))

	if injectFront {
		// label default container
		const annotationsKeyDefaultContainer = "kubectl.kubernetes.io/default-container"
		const annotationsKeyDefaultLogsContainer = "kubectl.kubernetes.io/default-logs-container"
		if pod.Labels == nil {
			pod.Labels = make(map[string]string, 1)
		}
		pod.Annotations[annotationsKeyDefaultContainer] = pod.Spec.Containers[0].Name
		pod.Annotations[annotationsKeyDefaultLogsContainer] = pod.Spec.Containers[0].Name
		pod.Spec.Containers = append([]corev1.Container{sidecarContainer}, pod.Spec.Containers...)
	} else {
		pod.Spec.Containers = append(pod.Spec.Containers, sidecarContainer)
	}

	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string, 2)
	}
	pod.Annotations["sidecar_injected_by"] = "webhook-kb-demo"
}
