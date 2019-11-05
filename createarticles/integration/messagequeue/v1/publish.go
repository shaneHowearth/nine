package messagequeue

type MQ interface {
	Publish(id string) error
}
