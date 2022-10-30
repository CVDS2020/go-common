package svc

func WithNotify(notify bool) Option {
	return optionFunc(func(service Service) {
		if systemdService, is := service.(*linuxSystemdService); is {
			systemdService.systemdNotify = notify
		}
	})
}
