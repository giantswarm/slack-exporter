package service

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"

	"github.com/giantswarm/slack-exporter/flag"
	"github.com/giantswarm/slack-exporter/service/collector"
)

type Config struct {
	Logger micrologger.Logger

	Description string
	Flag        *flag.Flag
	GitCommit   string
	ProjectName string
	Source      string
	Viper       *viper.Viper
}

type Service struct {
	Version *version.Service

	bootOnce          sync.Once
	exporterCollector *collector.Set
}

func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flag must not be empty", config)
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Viper must not be empty", config)
	}

	var err error

	var slackClient *slack.Client
	{
		slackClient = slack.New(config.Viper.GetString(config.Flag.Service.Slack.Auth.Token))
	}

	var exporterCollector *collector.Set
	{
		c := collector.SetConfig{
			Logger:      config.Logger,
			SlackClient: slackClient,

			ChannelExpressions: mustParseJSONList(config.Viper.GetString(config.Flag.Service.Collector.Channel.Expressions)),
		}

		exporterCollector, err = collector.NewSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		c := version.Config{
			Description: config.Description,
			GitCommit:   config.GitCommit,
			Name:        config.ProjectName,
			Source:      config.Source,
		}

		versionService, err = version.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Service{
		Version: versionService,

		bootOnce:          sync.Once{},
		exporterCollector: exporterCollector,
	}

	return s, nil
}

func (s *Service) Boot(ctx context.Context) {
	s.bootOnce.Do(func() {
		go s.exporterCollector.Boot(ctx)
	})
}

func mustParseJSONList(s string) []string {
	var l []string
	err := json.Unmarshal([]byte(s), &l)
	if err != nil {
		panic(err)
	}

	return l
}
