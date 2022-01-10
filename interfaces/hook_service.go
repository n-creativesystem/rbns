package interfaces

type LoginHook interface {
	Run()
}

type HooksService struct {
	loginHooks []LoginHook
}

func (s *HooksService) RunLoginHook() {
	for _, hook := range s.loginHooks {
		hook.Run()
	}
}
