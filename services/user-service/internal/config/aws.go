package config

type AWS struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	S3Bucket        string
	S3Endpoint      string
}

func loadAWS() (*AWS, error) {
	accessKeyID, err := getEnvRequired("AWS_ACCESS_KEY_ID")
	if err != nil {
		return nil, err
	}
	secretAccessKey, err := getEnvRequired("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		return nil, err
	}
	region, err := getEnvRequired("AWS_REGION")
	if err != nil {
		return nil, err
	}
	s3Bucket, err := getEnvRequired("AWS_S3_BUCKET")
	if err != nil {
		return nil, err
	}
	s3Endpoint, err := getEnvRequired("AWS_S3_ENDPOINT")
	if err != nil {
		return nil, err
	}

	return &AWS{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
		S3Bucket:        s3Bucket,
		S3Endpoint:      s3Endpoint,
	}, nil
}
