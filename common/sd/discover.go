package sd

type Discover interface {
	Endpoints(name string) ([]Service, error)
}
