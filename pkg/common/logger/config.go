package log

type Config struct {
	JSON  bool   `yaml:"json"`
	Level string `yaml:"level"`
	// true to enable log sampling, where the same log message and level will be throttled.
	// we have two layers of sampling
	// 1. global sampling - within a second, it will log the first SampleInitial, then every SampleInterval messages.
	// 2. per participant/track sampling - to be used with Logger.WithItemSampler(). This would be used to throttle
	//    the logs for a particular participant/track.
	Sample bool `yaml:"sample,omitempty"`

	// global sampling per server
	// when sampling, the first N logs will be logged
	SampleInitial int `yaml:"sample_initial,omitempty"`
	// when sampling, every Mth log will be logged
	SampleInterval int `yaml:"sample_interval,omitempty"`

	// participant/track level sampling
	ItemSampleSeconds  int `yaml:"item_sample_seconds,omitempty"`
	ItemSampleInitial  int `yaml:"item_sample_initial,omitempty"`
	ItemSampleInterval int `yaml:"item_sample_interval,omitempty"`
}
