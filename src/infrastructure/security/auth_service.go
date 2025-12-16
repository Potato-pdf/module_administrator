package security

import (
	"fmt"

	"luminaMO-orq-modAI/src/domain/services"
)

type authService struct {
	inboundKeys  map[string]string
	outboundKeys map[string]string
	permissions  map[string][]string
}

type AuthConfig struct {
	InboundKeys  map[string]string
	OutboundKeys map[string]string
	Permissions  map[string][]string
}

func NewAuthService(config AuthConfig) services.AuthService {
	return &authService{
		inboundKeys:  config.InboundKeys,
		outboundKeys: config.OutboundKeys,
		permissions:  config.Permissions,
	}
}

func (s *authService) ValidateAPIKey(apiKey string, expectedSource string) error {
	expectedKey, exists := s.inboundKeys[expectedSource]
	if !exists {
		return fmt.Errorf("unknown source: %s", expectedSource)
	}

	if apiKey != expectedKey {
		return fmt.Errorf("invalid API key for source: %s", expectedSource)
	}

	return nil
}

func (s *authService) GetAPIKeyForDestination(destination string) (string, error) {
	apiKey, exists := s.outboundKeys[destination]
	if !exists {
		return "", fmt.Errorf("no API key configured for destination: %s", destination)
	}

	return apiKey, nil
}

func (s *authService) ValidateModulePermission(apiKey string, module string) error {
	allowedModules, exists := s.permissions[apiKey]
	if !exists {
		return fmt.Errorf("API key has no permissions")
	}

	for _, allowedModule := range allowedModules {
		if allowedModule == module || allowedModule == "*" {
			return nil
		}
	}

	return fmt.Errorf("API key not authorized for module: %s", module)
}
