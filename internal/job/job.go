package job

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/murtaza-u/rinc/internal/conf"

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
	slog.SetDefault(logger(c.Log))
	return Job{
		conf:       c,
		kubeClient: kubeClient,
	}
}

// GenerateAll generates reports for all the configured tasks.
func (j Job) GenerateAll(ctx context.Context) error {
	now := time.Now().UTC()
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

func logger(c conf.Log) *slog.Logger {
	var level slog.Level
	switch c.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	opt := &slog.HandlerOptions{Level: level}
	if c.Format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stderr, opt))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, opt))
}
