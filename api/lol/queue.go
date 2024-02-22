package api

type QueueType int

const (
	normal QueueType = 400
	ranked QueueType = 420
	clash  QueueType = 700
)

func (q QueueType) String() string {
	s, _ := map[QueueType]string{
		normal: "Normal",
		ranked: "Ranked",
		clash:  "Clash",
	}[q]

	return s
}

type queues struct {
	Normal QueueType
	Ranked QueueType
	Clash  QueueType
}

var Queue = queues{
	Normal: normal,
	Ranked: ranked,
	Clash:  clash,
}
