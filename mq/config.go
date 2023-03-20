package mq

type Config struct {
	KafkaAddress string `json:"kafka_address" required:"true"`
	GroupName    string `json:"group_name"    required:"true"`
}
