package config

import "github.com/HomesNZ/go-common/env"

func New(region, queueName string) (*Config, error) {
	cfg := &Config{
		QueueName: queueName,
		Region:    region,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewFromEnv() (*Config, error) {
	region := env.GetString("AWS_SQS_REGION", "")
	queueName := env.GetString("AWS_SQS_QUEUE", "")
	maxMsg := env.GetInt("AWS_SQS_MAX_MESSAGES", 1)
	maxWorker := env.GetInt("AWS_SQS_MAX_WORKERS", 1)
	maxHandleTime := env.GetInt("AWS_SQS_MAX_HANDLE_TIME", 600) // 10 minutes
	maxNumHandlers := env.GetInt("AWS_SQS_MAX_HANDLERS", 10)

	var maxNumOfMessages int32
	if maxMsg > 100 {
		maxNumOfMessages = 100
	} else {
		maxNumOfMessages = int32(maxMsg)
	}

	var maxNumOfWorkers int
	if maxWorker > 100 {
		maxNumOfWorkers = 100
	} else {
		maxNumOfWorkers = maxWorker
	}

	var maxMessageHandleTime uint16
	if maxHandleTime > 43200 { // 12 hours
		maxMessageHandleTime = 43200
	} else {
		maxMessageHandleTime = uint16(maxHandleTime)
	}

	cfg := &Config{
		QueueName:            queueName,
		Region:               region,
		MaxMessageHandleTime: maxMessageHandleTime,
		MaxMsg:               maxNumOfMessages,
		MaxWorker:            maxNumOfWorkers,
		MaxHandlers:          maxNumHandlers,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
