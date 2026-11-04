package config

type URL struct {
	VerifyBaseURL         string
	ForgotPasswordBaseURL string
}

func loadURL() (*URL, error) {
	verifyBaseURL, err := getEnvRequired("VERIFY_BASE_URL")
	if err != nil {
		return nil, err
	}

	forgotPasswordBaseURL, err := getEnvRequired("FORGOT_PASSWORD_BASE_URL")
	if err != nil {
		return nil, err
	}

	return &URL{
		VerifyBaseURL:         verifyBaseURL,
		ForgotPasswordBaseURL: forgotPasswordBaseURL,
	}, nil
}
