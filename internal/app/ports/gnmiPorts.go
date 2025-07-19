package ports

type GNMIClient interface {
	Capabilities() (map[string]interface{}, error)
	GetAddress() string
}
