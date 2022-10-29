package eureka

// Eureka eureka registry
type Eureka struct {
	Server string `json:"server"`
	Url    string `json:"url"`
}

// NewEureka create a new eureka registry
func NewEureka(srv, url string) *Eureka {
	return &Eureka{
		Server: srv,
		Url:    url,
	}
}
