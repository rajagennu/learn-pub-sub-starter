package pubsub

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota + 1
	Transient
)

func (d SimpleQueueType) String() string {
	return [...]string{"durable", "transient", "South", "West"}[d-1]
}

func (d SimpleQueueType) EnumIndex() int {
	return int(d)
}
