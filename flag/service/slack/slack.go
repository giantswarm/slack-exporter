package slack

import (
	"github.com/giantswarm/slack-exporter/flag/service/slack/auth"
)

type Slack struct {
	Auth auth.Auth
}
