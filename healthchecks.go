package goinjection

type ApplicationHealthCheck interface {
	CheckHealth() error
}

func (this *Application) CheckHealth() error {
	for _, service := range this.services {
		if check, ok := interface{}(service).(ApplicationHealthCheck); ok {
			if err := check.CheckHealth(); err != nil {
				return err
			}

		}
	}
	return nil
}
