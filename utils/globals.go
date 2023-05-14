package utils

type SpaceType int

const (
	Ip     SpaceType = 1
	Cosine SpaceType = 2
	L2     SpaceType = 3
)

func SpaceFromString(space string) SpaceType {
	switch space {
	case "ip":
		return Ip
	case "cosine":
		return Cosine
	case "l2":
		return L2
	default:
		return Cosine
	}
}
