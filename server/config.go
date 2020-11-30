package server

type Config struct {
	Threshold  int
	Port       uint
	OnlineMode bool
}

func DefaultConfig() Config {
	return Config{
		Threshold:  255,
		Port:       25565,
		OnlineMode: false,
	}
}
