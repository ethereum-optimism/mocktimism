// Package orchestrator provides tools and utilities to enable service discovery using Zeroconf.
package orchestrator

import (
	"context"
	"log"

	"github.com/grandcat/zeroconf"
)

// Orchestrator manages service registration and discovery using Zeroconf.
type Orchestrator struct {
	resolver    *zeroconf.Resolver
	services    map[string]*zeroconf.ServiceEntry
	serviceType string
}

// Service represents the interface that a service should implement
// to be registered and discovered using Orchestrator.
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
	Config() interface{}
	// Starts the service.
	Start(ctx context.Context) error
}

// NewOrchestrator initializes and returns a new Orchestrator instance.
// The serviceType argument specifies the type of services that the instance will manage.
func NewOrchestrator(serviceType string) *Orchestrator {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalf("Failed to initialize resolver: %v", err)
	}

	return &Orchestrator{
		resolver:    resolver,
		services:    make(map[string]*zeroconf.ServiceEntry),
		serviceType: serviceType,
	}
}

// Register registers a given service with the Orchestrator.
// The provided service should implement the Service interface.
func (sd *Orchestrator) Register(s Service) {
	var txtRecords []string
	if configMap, ok := s.Config().(map[string]string); ok {
		txtRecords = make([]string, 0, len(configMap))
		for key, val := range configMap {
			txtRecords = append(txtRecords, key+"="+val)
		}
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

// GetServices returns a list of service IDs that are currently registered with the Orchestrator.
func (sd *Orchestrator) GetServices() []string {
	ids := make([]string, 0, len(sd.services))
	for id := range sd.services {
		ids = append(ids, id)
	}
	return ids
}

// GetServiceById retrieves a registered service based on its ID.
// Returns nil if the ID does not match any registered service.
func (sd *Orchestrator) GetServiceById(id string) *zeroconf.ServiceEntry {
	return sd.services[id]
}
