// Package servicediscovery provides tools and utilities to enable service discovery using Zeroconf.
package servicediscovery

import (
	"log"

	"github.com/grandcat/zeroconf"
)

// ServiceConfig is a type alias for a map containing service configuration key-value pairs.
type ServiceConfig map[string]string

// ServiceDiscovery manages service registration and discovery using Zeroconf.
type ServiceDiscovery struct {
	resolver    *zeroconf.Resolver
	services    map[string]*zeroconf.ServiceEntry
	serviceType string
}

// Service represents the interface that a service should implement
// to be registered and discovered using ServiceDiscovery.
type Service interface {
	// Returns the host name of the service.
	Hostname() string
	// Returns the port number on which the service is listening.
	Port() int
	// Returns the type of the service, e.g., "_myService._tcp".
	ServiceType() string
	// Returns a unique identifier for the service.
	ID() string
	// Returns a map containing service configuration key-value pairs.
	Config() ServiceConfig
	// Starts the service.
	Start() error
}

// NewServiceDiscovery initializes and returns a new ServiceDiscovery instance.
// The serviceType argument specifies the type of services that the instance will manage.
func NewServiceDiscovery(serviceType string) *ServiceDiscovery {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalf("Failed to initialize resolver: %v", err)
	}

	return &ServiceDiscovery{
		resolver:    resolver,
		services:    make(map[string]*zeroconf.ServiceEntry),
		serviceType: serviceType,
	}
}

// Register registers a given service with the ServiceDiscovery.
// The provided service should implement the Service interface.
func (sd *ServiceDiscovery) Register(s Service) {
	txtRecords := make([]string, 0, len(s.Config()))
	for key, val := range s.Config() {
		txtRecords = append(txtRecords, key+"="+val)
	}

	server, err := zeroconf.Register(s.Hostname(), s.ServiceType(), "local.", s.Port(), txtRecords, nil)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Store service for future reference
	sd.services[s.ID()] = &zeroconf.ServiceEntry{
		HostName: s.Hostname(),
		Port:     s.Port(),
		Text:     txtRecords,
	}

	defer server.Shutdown()
}

// GetServices returns a list of service IDs that are currently registered with the ServiceDiscovery.
func (sd *ServiceDiscovery) GetServices() []string {
	ids := make([]string, 0, len(sd.services))
	for id := range sd.services {
		ids = append(ids, id)
	}
	return ids
}

// GetServiceById retrieves a registered service based on its ID.
// Returns nil if the ID does not match any registered service.
func (sd *ServiceDiscovery) GetServiceById(id string) *zeroconf.ServiceEntry {
	return sd.services[id]
}
