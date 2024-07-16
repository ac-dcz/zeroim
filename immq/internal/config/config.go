package config

type Config struct {
	Kq struct {
		Brokers []string
		Topic   string
		GroupID string
	}
}
