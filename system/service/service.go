package svc

// Service interface specifies the behavior of the system service
type Service interface {
	// Run method is the main function of program running. The Run method
	// needs to include the implementation of notifying the system service
	// daemon(e.g. linux systemd or Windows service manager) when the
	// program starts, exits, or has errors
	Run() (exitCode int)
}
