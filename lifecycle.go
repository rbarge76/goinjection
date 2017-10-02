package goinjection

type ApplicationSetup interface {
	DoSetup() error
}

type ApplicationShutdown interface {
	Shutdown()
}

