package internal

type Service int

const (
	NoneService Service = iota
	QueryService
	IndexService
	StoreService
)

var serviceNames = [...]string{
	"NoneService",
	"QueryService",
	"IndexService",
	"StoreService",
}

func (s Service) String() string {
	if s < NoneService || s > StoreService {
		return "NoneService"
	}
	return serviceNames[s]
}
