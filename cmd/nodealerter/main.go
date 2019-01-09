package main

import (
	"context"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/richardcase/nodealerter/pkg/controller"
	"github.com/richardcase/nodealerter/pkg/k8s"
	"github.com/richardcase/nodealerter/pkg/signal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var (
	options *Options

	rootCmd = &cobra.Command{
		Use:   "nodealerter",
		Short: "Node Alerter will alert if the number of pods exceeds a threshold",
		Long:  "",
		Run: func(c *cobra.Command, _ []string) {
			if err := doRun(); err != nil {
				logrus.Fatalf("error running node alerter: %v", err)
			}
		},
	}
)

func main() {
	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal("error executing root command")
		os.Exit(1)
	}
}

func init() {
	options = &Options{}

	//NOTE: this would really be a command line argument
	logrus.SetLevel(logrus.DebugLevel)

	rootCmd.PersistentFlags().StringVar(&options.ConfigFile,
		"config", "",
		"Config file (default is $HOME/.nodealerter.yaml)")
	rootCmd.PersistentFlags().StringVarP(&options.KubeconfigFile,
		"kubeconfig", "k", "",
		"Absolute path to the kubeconfig file. Only required if out-of-cluster.")
	rootCmd.PersistentFlags().IntVarP(&options.NodeThreshold, "nodes", "", 10, "the number of nodes after which an alert will be raised. Defaults to 10.")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	_ = viper.BindPFlag("nodes", rootCmd.PersistentFlags().Lookup("nodes"))
}

func initConfig() {
	if options.ConfigFile != "" {
		viper.SetConfigFile(options.ConfigFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			logrus.Fatalf("failed to get home directory: %v", err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".nodealerter")
	}

	replacer := strings.NewReplacer(".", "-")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Debug("using config file: ", viper.ConfigFileUsed())
	}
}

func doRun() error {
	stopChan := signal.SetupSignalHandler()

	kubeConfig := k8s.MustGetClientConfig(options.KubeconfigFile)
	kubeClient := k8s.MustGetKubeClientFromConfig(kubeConfig)

	controllerConfig := &controller.Config{
		KubeConfig:     kubeConfig,
		KubeClient:     kubeClient,
		NodesThreshold: options.NodeThreshold,
	}

	ct := controller.New(*controllerConfig)

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error { return ct.Run(ctx) })

	logrus.Info("started node alerter controller")

	select {
	case <-stopChan:
		logrus.Info(("shutdown signal received, shutdown...."))
	case <-ctx.Done():
	}

	cancel()
	if err := wg.Wait(); err != nil {
		return errors.Wrap(err, "unhandled error, exiting")
	}
	return nil
}
