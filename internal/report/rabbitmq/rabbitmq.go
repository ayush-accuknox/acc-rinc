package rabbitmq

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/murtaza-u/rinc/internal/conf"
	"github.com/murtaza-u/rinc/internal/util"
	"github.com/murtaza-u/rinc/view/layout"
	"github.com/murtaza-u/rinc/view/partial"
	tmpl "github.com/murtaza-u/rinc/view/rabbitmq"

	"k8s.io/client-go/kubernetes"
)

// Reporter is the rabbitmq health metrics reporter.
type Reporter struct {
	kubeClient *kubernetes.Clientset
	conf       conf.RabbitMQ
}

// NewReporter creates a new of the rabbitmq reporter.
func NewReporter(c conf.RabbitMQ, kubeClient *kubernetes.Clientset) Reporter {
	return Reporter{
		conf:       c,
		kubeClient: kubeClient,
	}
}

// Report satisfies the report.Reporter interface by writing the RabbitMQ
// cluster status and fetched metrics to the provided io.Writer.
func (r Reporter) Report(ctx context.Context, to io.Writer, now time.Time) error {
	stamp := now.Format(util.IsosecLayout)
	up, err := r.IsClusterUp(ctx)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"fetching rabbitmq health status",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("fetching rabbitmq health status: %w", err)
	}
	if !up {
		slog.LogAttrs(
			ctx,
			slog.LevelInfo,
			"rabbitmq cluster is down",
		)
		c := layout.Base(
			fmt.Sprintf("RabbitMQ - %s | AccuKnox Reports", stamp),
			partial.Navbar(false, false),
			tmpl.Report(tmpl.Data{
				Timestamp: now,
				IsHealthy: false,
			}),
		)
		err := c.Render(ctx, to)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"rendering rabbitmq template",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("rendering rabbitmq template: %w", err)
		}
		return nil
	}
	metrics, err := r.GetMetrics(ctx)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to fetch rabbitmq metrics",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to fetch rabbitmq metrics: %w", err)
	}
	c := layout.Base(
		fmt.Sprintf("RabbitMQ - %s | AccuKnox Reports", stamp),
		partial.Navbar(false, false),
		tmpl.Report(tmpl.Data{
			Timestamp: now,
			IsHealthy: true,
			Metrics:   *metrics,
		}),
	)
	err = c.Render(ctx, to)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"rendering rabbitmq template",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("rendering rabbitmq template: %w", err)
	}
	return nil
}
