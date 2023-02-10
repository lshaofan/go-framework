package micro

type Option func(o interface{})
type Service interface {
	Run() error
	Close()
	Init() error
}
