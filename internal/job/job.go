package job

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/murtaza-u/rinc/internal/conf"
	"github.com/murtaza-u/rinc/internal/report/rabbitmq"
	"github.com/murtaza-u/rinc/internal/util"
	"github.com/murtaza-u/rinc/view"
	"github.com/murtaza-u/rinc/view/layout"
	"github.com/murtaza-u/rinc/view/partial"

	"k8s.io/client-go/kubernetes"
)

// Job runs inside the Kubernetes cluster and generates status and metrics
// reports.
type Job struct {
	conf       conf.C
	kubeClient *kubernetes.Clientset
}

// New returns a new reporting Job object.
func New(c conf.C, kubeClient *kubernetes.Clientset) Job {
	slog.SetDefault(util.NewLogger(c.Log))
	return Job{
		conf:       c,
		kubeClient: kubeClient,
	}
}

// GenerateAll generates reports for all the configured tasks.
func (j Job) GenerateAll(ctx context.Context) error {
	now := time.Now().UTC()

	err := j.GenerateIndex(ctx, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to generate index page",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to generate index page: %w", err)
	}

	if j.conf.RabbitMQ.Enable {
		err := j.GenerateRMQReport(ctx, now)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"generating RMQ report",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("generating RMQ report: %w", err)
		}
	}

	return nil
}

// GenerateIndex generates the index status page.
func (j Job) GenerateIndex(ctx context.Context, now time.Time) error {
	stamp := now.Format(util.IsosecLayout)
	var statuses []view.IndexStatus

	if j.conf.RabbitMQ.Enable {
		rmqR := rabbitmq.NewReporter(j.conf.RabbitMQ, j.kubeClient)
		up, err := rmqR.IsClusterUp(ctx)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"checking rabbitmq health status",
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("checking rabbitmq health: %w", err)
		}
		statuses = append(statuses, view.IndexStatus{
			Name:    "RabbitMQ",
			ID:      stamp,
			Healthy: up,
		})
	}

	if len(statuses) == 0 {
		return nil
	}

	err := j.initStamp(ctx, stamp)
	if err != nil {
		return fmt.Errorf("initializing %q stamp: %w", stamp, err)
	}

	path := filepath.Join(j.conf.Output, stamp, "index.html")
	f, err := os.Create(path)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			fmt.Sprintf("creating %q file", path),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("creating %q file: %w", path, err)
	}
	defer f.Close()

	c := layout.Base(
		fmt.Sprintf("%s | AccuKnox Reports", stamp),
		partial.Navbar(false, false),
		view.Index(statuses),
		partial.Footer(now),
	)
	err = c.Render(ctx, f)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"rendering index.html",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("rendering index.html: %w", err)
	}

	return nil
}

func (j Job) initStamp(ctx context.Context, stamp string) error {
	path := filepath.Join(j.conf.Output, stamp)
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			fmt.Sprintf("creating %q directory", path),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("creating %q directory: %w", path, err)
	}
	return nil
}
