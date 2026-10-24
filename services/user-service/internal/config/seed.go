package config

type Seed struct {
	Email    string
	Password string
}

func loadSeed() (*Seed, error) {
	email, err := getEnvRequired("SEED_ADMIN_EMAIL")
	if err != nil {
		return nil, err
	}

	password, err := getEnvRequired("SEED_ADMIN_PASSWORD")
	if err != nil {
		return nil, err
	}

	return &Seed{
		Email:    email,
		Password: password,
	}, nil
}
