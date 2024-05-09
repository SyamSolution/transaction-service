package config

type BaseDep struct {
	Logger Logger
}

func NewBaseDep() *BaseDep {
	return &BaseDep{
		Logger: SetupLogger(),
	}
}
