package service

import (
	"github.com/giantswarm/slack-exporter/flag/service/collector"
	"github.com/giantswarm/slack-exporter/flag/service/slack"
)

type Service struct {
	Collector collector.Collector
	Slack     slack.Slack
}
