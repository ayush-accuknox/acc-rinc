package main

import (
	"context"
	"log"
	"os"

	"github.com/murtaza-u/rinc/internal/conf"
	"github.com/murtaza-u/rinc/internal/job"
	"github.com/murtaza-u/rinc/internal/kube"
	"github.com/murtaza-u/rinc/internal/web"
)

func main() {
	conf, err := conf.New(os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	err = conf.Validate()
	if err != nil {
		log.Fatalf("validating provided config: %s", err.Error())
	}

	if conf.RunAsGenerator {
		kubeClient, err := kube.NewClient(conf.KubernetesClient)
		if err != nil {
			log.Fatalf("kubernetes client: %s", err.Error())
		}
		job := job.New(*conf, kubeClient)
		err = job.GenerateAll(context.Background())
		if err != nil {
			log.Fatalf("generating reports: %s", err.Error())
		}
		return
	}

	web.NewSrv(*conf).Run(context.Background())
}
