package application

import (
	"context"

	"gosuda.org/boilerplate/internal/domain"
	"gosuda.org/boilerplate/internal/infrastructure"
)

// DebugService handles debug-related operations
type DebugService struct {
	logger  infrastructure.LoggerInterface
	store   domain.Store
	metrics *infrastructure.Metrics
}

// NewDebugService creates a new debug service
func NewDebugService(logger infrastructure.LoggerInterface, store domain.Store, metrics *infrastructure.Metrics) *DebugService {
	return &DebugService{
		logger:  logger,
		store:   store,
		metrics: metrics,
	}
}

// SetLogLevel changes the application log level at runtime
func (s *DebugService) SetLogLevel(ctx context.Context, level string) error {
	// Validate log level
	if err := validateLogLevel(level); err != nil {
		return err
	}

	// Set the log level
	if err := s.logger.SetLevel(level); err != nil {
		return &domain.StorageError{Err: err}
	}

	// Log the level change
	s.logger.LogLevelChange(s.logger.GetLevel().String(), level)

	return nil
}

// GetLogLevel returns the current log level
func (s *DebugService) GetLogLevel(ctx context.Context) string {
	return s.logger.GetLevel().String()
}

// GetMetrics returns application metrics
func (s *DebugService) GetMetrics(ctx context.Context) (string, error) {
	return s.metrics.Gather()
}

// GetHealthStatus returns the application health status
func (s *DebugService) GetHealthStatus(ctx context.Context) (*HealthStatus, error) {
	// Check storage health
	storageHealthy := true
	if err := s.checkStorageHealth(ctx); err != nil {
		storageHealthy = false
	}

	// Check logger health
	loggerHealthy := true
	// Logger is always available since it's required in constructor

	status := &HealthStatus{
		Status: "healthy",
		Checks: map[string]HealthCheck{
			"storage": {
				Status:  storageHealthy,
				Message: "Storage is accessible",
			},
			"logger": {
				Status:  loggerHealthy,
				Message: "Logger is available",
			},
		},
	}

	// Determine overall status
	if !storageHealthy || !loggerHealthy {
		status.Status = "unhealthy"
	}

	return status, nil
}

// HealthStatus represents the application health status
type HealthStatus struct {
	Status string                 `json:"status"`
	Checks map[string]HealthCheck `json:"checks,omitempty"`
}

// HealthCheck represents a health check result
type HealthCheck struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// validateLogLevel validates a log level string
func validateLogLevel(level string) error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[level] {
		return &domain.ValidationError{
			Field:   "level",
			Message: "invalid log level",
		}
	}

	return nil
}

// checkStorageHealth checks if the storage is healthy
func (s *DebugService) checkStorageHealth(ctx context.Context) error {
	// Try to perform a simple operation to check storage health
	testKey := "health:test"
	testValue := "test"

	if err := s.store.Set(testKey, testValue); err != nil {
		return err
	}

	if _, err := s.store.Get(testKey); err != nil {
		return err
	}

	if err := s.store.Delete(testKey); err != nil {
		return err
	}

	return nil
}
