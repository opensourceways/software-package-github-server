package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/opensourceways/robot-github-lib/client"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/opensourceways/server-common-lib/secret"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/software-package-github-server/mq"
)

type options struct {
	service liboptions.ServiceOptions
	github  liboptions.GithubOptions
}

func (o *options) Validate() error {
	if err := o.service.Validate(); err != nil {
		return err
	}

	return o.github.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.github.AddFlags(fs)
	o.service.AddFlags(fs)

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit("software-package")
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options: %v", err)
	}

	cfg, err := LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Fatalf("load config file failed: %v", err)
	}

	if err = mq.Init(&cfg.MQ, log); err != nil {
		logrus.Fatalf("initialize mq failed, err:%v", err)
	}

	defer mq.Exit()

	secretAgent := new(secret.Agent)
	if err = secretAgent.Start([]string{o.github.TokenPath}); err != nil {
		logrus.Fatalf("starting secret agent error: %v", err)
	}

	defer secretAgent.Stop()

	c := client.NewClient(secretAgent.GetTokenGenerator(o.github.TokenPath))

	m := newMessageService(c, cfg)
	if err = m.subscribe(); err != nil {
		logrus.Fatalf("subscribe failed, err:%v", err)
	}

	run()
}

func run() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			logrus.Info("receive done. exit normally")

			return

		case <-sig:
			logrus.Info("receive exit signal")
			done()
			called = true

			return
		}
	}(ctx)

	<-ctx.Done()
}
