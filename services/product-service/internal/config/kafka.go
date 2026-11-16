package config

import "strings"

type Kafka struct {
	Brokers []string
}

func LoadKafka() (*Kafka, error) {
	brokersValue, err := getEnvRequired("KAFKA_BROKERS")
	if err != nil {
		return nil, err
	}

	brokers := make([]string, 0)

	for _, broker := range strings.Split(brokersValue, ",") {
		broker = strings.TrimSpace(broker)

		if broker == "" {
			continue
		}
		brokers = append(brokers, broker)
	}

	return &Kafka{
		Brokers: brokers,
	}, nil
}
