package redis

func New(opts ...DeclareConfig) (Cache, error) {
	cfg := &Config{}
	if opts != nil {
		for _, opt := range opts {
			opt(cfg)
		}
	}
	if err := cfg.Validate(); err != nil {
		cfg, err = ConfigFromEnv()
		if err != nil {
			return nil, err
		}
	}

	return NewCache(cfg)
}
