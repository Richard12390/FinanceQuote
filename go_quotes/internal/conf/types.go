package conf

type Pipeline struct {
	Name                string   `yaml:"name"`
	Source              string   `yaml:"source"`
	AssetType           string   `yaml:"type"`
	Symbols             []string `yaml:"symbols"`
	Every               string   `yaml:"every"`
	BatchSize           int      `yaml:"batch_size"`
	BatchesPerTick      int      `yaml:"batches_per_tick"`
	MaxWorkers          int      `yaml:"max_workers"`
	MaxConcurrency      int      `yaml:"max_concurrency"`
	OutputDir           string   `yaml:"output_dir"`
	WriteJson           *bool    `yaml:"write_json,omitempty"`
	WriteJsonEveryTicks int      `yaml:"write_json_every_ticks,omitempty"`
}

type Config struct {
	NATSUrl             string     `yaml:"nats_url,omitempty"`
	EnablePersistAck    bool       `yaml:"enable_persist_ack,omitempty"`
	WriteJson           bool       `yaml:"write_json,omitempty"`
	WriteJsonEveryTicks int        `yaml:"write_json_every_ticks,omitempty"`
	Pipelines           []Pipeline `yaml:"pipelines,omitempty"`
}
