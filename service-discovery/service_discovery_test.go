package servicediscovery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MockService struct{}

func (m *MockService) Hostname() string {
	return "MyHost"
}

func (m *MockService) Port() int {
	return 8080
}

func (m *MockService) ServiceType() string {
	return "_myService._tcp"
}

func (m *MockService) ID() string {
	return "service123"
}

func (m *MockService) Config() ServiceConfig {
	return ServiceConfig{
		"version": "1.0",
		"env":     "production",
	}
}

func (m *MockService) Start() error {
	// As it's a mock service, you might not have any specific logic to run.
	// But if you had an HTTP server or some logic here, you'd start it.
	return nil
}

func TestServiceDiscovery(t *testing.T) {
	sd := NewServiceDiscovery("_workstation._tcp")

	sd.Register(&MockService{})

	// 3. Get all services
	allServices := sd.GetServices()
	require.Contains(t, allServices, "service123", "Service ID not found in the registered services list")

	// 4. Fetch service by ID
	service := sd.GetServiceById("service123")
	require.NotNil(t, service, "Expected service not to be nil")

	// 5. Validate service details
	require.Equal(t, "MyHost", service.HostName)
	require.Equal(t, 8080, service.Port)
	require.Contains(t, service.Text, "version=1.0")
	require.Contains(t, service.Text, "env=production")
}
