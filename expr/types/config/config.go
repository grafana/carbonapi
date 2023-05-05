package config

type ConfigType = struct {
	// NudgeStartTimeOnAggregation enables nudging the start time of metrics when aggregated
	// to honor MaxDataPoints. The start time is nudged in such way that timestamps always
	// fall in the same bucket. This is done bY GraphiteWeb, and is useful to avoid jitter
	// in graphs when refreshing the page.
	NudgeStartTimeOnAggregation bool

	// TODO
	UseBucketsHighestTimestampOnAggregation bool
}

var Config = ConfigType{}
