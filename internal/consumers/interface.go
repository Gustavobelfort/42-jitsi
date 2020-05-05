package consumers

// Consumer defines a model meant to read from a scale_team source and insert it with its passed handler.
//
// Different types of consumers will be implemented, like a direct API endpoint or a RabbitMQ consumer.
type Consumer interface {
	Start() error
	Stop() error
}
