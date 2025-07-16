package ports

type GNMIClient interface {
	Connect() error
	Capabilities() (map[string][]string, error)
	Close() error
	Target() string
}
