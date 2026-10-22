package config

const (
	defaultAppEnv         = "development"
	defaultAppPort        = "8001"
	defaultAppIsDebugMode = false
)

type App struct {
	Env         string
	Port        string
	IsDebugMode bool
}

func loadApp() App {
	return App{
		Env:         getEnv("APP_ENV", defaultAppEnv),
		Port:        getEnv("APP_PORT", defaultAppPort),
		IsDebugMode: getEnvBool("APP_IS_DEBUG_MODE", defaultAppIsDebugMode),
	}
}

func (a *App) validate() error {
	if err := validatePort("APP_PORT", a.Port); err != nil {
		return err
	}
	return nil
}
