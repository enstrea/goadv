package srv

type Server interface {
	Name() string
	Start() error
	Stop() error
}
