package module

type sink struct {
}

func (s *sink) OnStart() {
}

func (s *sink) OnTick() {
	defaultModuleMgr.onTick()
}

func (s *sink) OnStop() {
}


