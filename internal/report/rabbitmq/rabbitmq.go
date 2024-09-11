package rabbitmq

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/murtaza-u/rinc/internal/conf"
	"github.com/murtaza-u/rinc/internal/report"

	"k8s.io/client-go/kubernetes"
)

// Reporter is the rabbitmq health metrics reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.RabbitMQ
}

// NewReporter creates a new of the rabbitmq reporter.
func NewReporter(c conf.RabbitMQ, kubeClient *kubernetes.Clientset) report.Reporter {
	return Reporter{
		conf:       c,
		kubeClient: kubeClient,
	}
}

// Report satisfies the report.Reporter interface by writing the RabbitMQ
// cluster status and fetched metrics to the provided io.Writer.
func (r Reporter) Report(ctx context.Context, to io.Writer) error {
	up, err := r.IsClusterUp(ctx)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to fetch rabbitmq health status",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to fetch rabbitmq health status: %w", err)
	}
	if !up {
		slog.LogAttrs(
			ctx,
			slog.LevelInfo,
			"rabbitmq cluster is down",
		)
		// TODO: render template
		return nil
	}
	_, err = r.GetMetrics(ctx)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to fetch rabbitmq metrics",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to fetch rabbitmq metrics: %w", err)
	}
	// TODO: render template
	return nil
}
