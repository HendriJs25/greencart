package config

type Redis struct {
	URL string
}

func loadRedis() (*Redis, error) {
	URL, err := getEnvRequired("REDIS_URL")
	if err != nil {
		return nil, err
	}

	return &Redis{URL: URL}, nil
}
