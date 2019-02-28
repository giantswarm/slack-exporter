package collector

import (
	"context"
	"fmt"
	"regexp"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/nlopes/slack"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	channelMembersDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "members_count"),
		"Slack channel members.",
		[]string{
			labelOrg,
			labelRepo,
			labelChannel,
		},
		nil,
	)
)

type IssueConfig struct {
	Logger      micrologger.Logger
	SlackClient *slack.Client

	ChannelExpressions []string
}

type Issue struct {
	logger      micrologger.Logger
	slackClient *slack.Client

	channelExpressions []string
}

func NewIssue(config IssueConfig) (*Issue, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.SlackClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.SlackClient must not be empty", config)
	}

	i := &Issue{
		logger:      config.Logger,
		slackClient: config.SlackClient,

		channelExpressions: config.ChannelExpressions,
	}

	return i, nil
}

func (i *Issue) Collect(ch chan<- prometheus.Metric) error {
	ctx := context.Background()
	channelMembers := map[string]float64{}

	cursor := ""
	page := 1
	for {
		p := &slack.GetConversationsParameters{
			Cursor: cursor,
		}
		l, c, err := i.slackClient.GetConversations(p)
		if err != nil {
			return microerror.Mask(err)
		}

		i.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("collecting %3d channels of page %2d", len(l), page))

		for _, channel := range l {
			if !matchesExpression(channel.Name, i.channelExpressions) {
				continue
			}

			channelMembers[channel.Name] = float64(channel.NumMembers)
		}

		cursor = c
		page++

		if cursor == "" {
			i.logger.LogCtx(ctx, "level", "debug", "message", "collected all channels")
			break
		}
	}

	for k, v := range channelMembers {
		ch <- prometheus.MustNewConstMetric(
			channelMembersDesc,
			prometheus.GaugeValue,
			v,
			githubOrg,
			githubRepo,
			k,
		)
	}

	return nil
}

func (i *Issue) Describe(ch chan<- *prometheus.Desc) error {
	ch <- channelMembersDesc
	return nil
}

func matchesExpression(channel string, expressions []string) bool {
	for _, e := range expressions {
		if regexp.MustCompile(e).MatchString(channel) {
			return true
		}
	}

	return false
}
