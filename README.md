[![CircleCI](https://circleci.com/gh/giantswarm/slack-exporter.svg?&style=shield)](https://circleci.com/gh/giantswarm/slack-exporter)

# slack-exporter

The slack-exporter exports Prometheus metrics for Slack data.



### Example Execution

```
./slack-exporter daemon --service.collector.channel.expressions='[ "^news-.*", "^sig-.*", "^support-.*" ]' --service.slack.auth.token=$(cat ~/.credential/slack-exporter-slack-token)
```



### Example Queries

Showing a graph of the number of members per channel.

```
slack_exporter_channel_members_count
```

Showing a graph of the number of members in support channels.

```
slack_exporter_channel_members_count{channel=~"support-.*"}
```
