package internal

type Service struct {
	Code uint32
	Name string
}

var (
	NoneService  = Service{Code: 0, Name: "NoneService"}
	QueryService = Service{Code: 1, Name: "QueryService"}
	IndexService = Service{Code: 2, Name: "IndexService"}
	StoreService = Service{Code: 3, Name: "StoreService"}
)

func (s Service) String() string {
	return s.Name
}
