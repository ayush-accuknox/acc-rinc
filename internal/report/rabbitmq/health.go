package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/accuknox/rinc/internal/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IsClusterUp checks whether the RabbitMQ cluster is running by first
// verifying that at least one RabbitMQ pod is in the READY state, followed by
// calling the management health check endpoint.
func (r Reporter) IsClusterUp(ctx context.Context) (bool, error) {
	ips, err := net.LookupIP(r.conf.HeadlessSvcAddr)
	if err != nil {
		return false, fmt.Errorf("lookup %q: %w", r.conf.HeadlessSvcAddr, err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"rabbitmq node ips",
		slog.Any("ips", ips),
	)

	ns, err := util.GetNamespaceFromFQDN(r.conf.HeadlessSvcAddr)
	if err != nil {
		return false, fmt.Errorf("failed to get namespace from %q: %w",
			r.conf.HeadlessSvcAddr, err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"rabbitmq namespace",
		slog.String("namespace", ns),
	)

	pods, err := r.kubeClient.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("list pods in %q namespace: %w", ns, err)
	}

	atleastOneUp := false
	for _, pod := range pods.Items {
		isRabbit := false
		for _, ip := range ips {
			if ip.String() == pod.Status.PodIP {
				isRabbit = true
				break
			}
		}
		if !isRabbit {
			continue
		}

		isReady := false
		for _, condition := range pod.Status.Conditions {
			if condition.Type != corev1.PodReady {
				continue
			}
			if condition.Status == corev1.ConditionTrue {
				isReady = true
			}
			break
		}

		status := pod.Status.Phase
		if isReady {
			atleastOneUp = true
			slog.LogAttrs(
				ctx,
				slog.LevelInfo,
				"at least one rabbitmq node is up",
				slog.String("name", pod.GetName()),
				slog.String("status", string(status)),
			)
			break
		}
		slog.LogAttrs(
			ctx,
			slog.LevelInfo,
			"rabbitmq node down",
			slog.String("name", pod.GetName()),
			slog.String("status", string(status)),
		)
	}

	if !atleastOneUp {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			fmt.Sprintf("no rabbitmq node ready in namespace %q", ns),
		)
		return false, nil
	}

	status, err := r.callEndpointReturnStatus(ctx, healthCheckEndpoint)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelInfo,
			"calling rabbitmq health check endpoint",
			slog.String("error", err.Error()),
		)
	}
	if status != 200 {
		return false, nil
	}

	return true, nil
}
