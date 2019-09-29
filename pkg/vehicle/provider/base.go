package provider

// Base represents the basic configuration for a cloud provider
type Base struct {
	Provider string `json:"provider"`
	UUID     string `json:"uuid"`
}

// GetProvider returns the provider's provider
func (c Base) GetProvider() string {
	return c.Provider
}
