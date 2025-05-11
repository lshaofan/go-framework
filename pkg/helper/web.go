package helper

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

// Run runs your Service.
//
// Run will block until one of the signals specified in sig is received or a provided context is done.
// If sig is empty syscall.SIGINT and syscall.SIGTERM are used by default.
func Run(service Service, sig ...os.Signal) error {
	env := environment{}
	if err := service.Init(env); err != nil {
		return err
	}

	if err := service.Start(); err != nil {
		return err
	}

	if len(sig) == 0 {
		sig = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	signalChan := make(chan os.Signal, 1)
	signalNotify(signalChan, sig...)

	var ctx context.Context
	if s, ok := service.(Context); ok {
		ctx = s.Context()
	} else {
		ctx = context.Background()
	}

	for {
		select {
		case s := <-signalChan:
			if h, ok := service.(Handler); ok {
				if err := h.Handle(s); err == ErrStop {
					goto stop
				}
			} else {
				// this maintains backwards compatibility for Services that do not implement Handle()
				goto stop
			}
		case <-ctx.Done():
			goto stop
		}
	}

stop:
	return service.Stop()

}

type environment struct{}

func (environment) IsWindowsService() bool {
	return false
}

// Create variable signal.Notify function so we can mock it in tests
var signalNotify = signal.Notify

var ErrStop = errors.New("stopping service")

// Service interface contains Start and Stop methods which are called
// when the service is started and stopped. The Init method is called
// before the service is started, and after it's determined if the program
// is running as a Windows Service.
//
// The Start and Init methods must be non-blocking.
//
// Implement this interface and pass it to the Run function to start your program.
type Service interface {
	// Init is called before the program/service is started and after it's
	// determined if the program is running as a Windows Service. This method must
	// be non-blocking.
	Init(Environment) error

	// Start is called after Init. This method must be non-blocking.
	Start() error

	// Stop is called in response to syscall.SIGINT, syscall.SIGTERM, or when a
	// Windows Service is stopped.
	Stop() error
}

// Context interface contains an optional Context function which a Service can implement.
// When implemented the context.Done() channel will be used in addition to signal handling
// to exit a process.
type Context interface {
	Context() context.Context
}

// Environment contains information about the environment
// your application is running in.
type Environment interface {
	// IsWindowsService reports whether the program is running as a Windows Service.
	IsWindowsService() bool
}

// Handler is an optional interface a Service can implement.
// When implemented, Handle() is called when a signal is received.
// Returning ErrStop from this method will result in Service.Stop() being called.
type Handler interface {
	Handle(os.Signal) error
}
