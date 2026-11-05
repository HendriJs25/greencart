package config

type AWS struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	S3Bucket        string
	S3Endpoint      string
	S3PublicBaseURL string
}

func loadAWS() (*AWS, error) {
	region, err := getEnvRequired("AWS_REGION")
	if err != nil {
		return nil, err
	}
	s3Bucket, err := getEnvRequired("AWS_S3_BUCKET")
	if err != nil {
		return nil, err
	}
	s3PublicBaseURL, err := getEnvRequired("AWS_S3_PUBLIC_BASE_URL")
	if err != nil {
		return nil, err
	}

	return &AWS{
		AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		Region:          region,
		S3Bucket:        s3Bucket,
		S3Endpoint:      getEnv("AWS_S3_ENDPOINT", ""),
		S3PublicBaseURL: s3PublicBaseURL,
	}, nil
}
