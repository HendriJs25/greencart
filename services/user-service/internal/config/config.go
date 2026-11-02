package config

type Config struct {
	App           App
	Database      *Database
	Seed          *Seed
	JWT           *JWT
	Redis         *Redis
	Kafka         *Kafka
	VerifyBaseURL string
}

func Load() (*Config, error) {
	database, err := loadDatabase()
	if err != nil {
		return nil, err
	}

	seed, err := loadSeed()
	if err != nil {
		return nil, err
	}

	jwt, err := loadJWT()
	if err != nil {
		return nil, err
	}

	redis, err := loadRedis()
	if err != nil {
		return nil, err
	}

	kafka, err := LoadKafka()
	if err != nil {
		return nil, err
	}

	verifyBaseURL, err := getEnvRequired("VERIFY_BASE_URL")
	if err != nil {
		return nil, err
	}

	return &Config{
		App:           loadApp(),
		Database:      database,
		Seed:          seed,
		JWT:           jwt,
		Redis:         redis,
		Kafka:         kafka,
		VerifyBaseURL: verifyBaseURL,
	}, nil
}

func (c *Config) validate() error {
	if err := c.App.validate(); err != nil {
		return err
	}

	if err := c.Database.validate(); err != nil {
		return err
	}

	if err := c.JWT.validate(); err != nil {
		return err
	}

	return nil
}

func (c *Config) ServerAddress() string {
	return ":" + c.App.Port
}
