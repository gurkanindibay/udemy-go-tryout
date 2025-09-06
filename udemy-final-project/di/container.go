package di

import (
	"github.com/gurkanindibay/udemy-rest-api/services"
	"github.com/samber/do/v2"
)

// Container holds all dependencies
type Container struct {
	Injector *do.RootScope
}

// NewContainer creates a new DI container
func NewContainer() *Container {
	injector := do.New()

	// Register services
	do.ProvideNamedValue(injector, "userService", services.NewUserService())
	do.ProvideNamedValue(injector, "eventService", services.NewEventService())
	do.ProvideNamedValue(injector, "authService", services.NewAuthService())

	return &Container{
		Injector: injector,
	}
}

// GetUserService returns the user service from the container
func (c *Container) GetUserService() services.UserService {
	return do.MustInvokeNamed[services.UserService](c.Injector, "userService")
}

// GetEventService returns the event service from the container
func (c *Container) GetEventService() services.EventService {
	return do.MustInvokeNamed[services.EventService](c.Injector, "eventService")
}

// GetAuthService returns the auth service from the container
func (c *Container) GetAuthService() services.AuthService {
	return do.MustInvokeNamed[services.AuthService](c.Injector, "authService")
}
