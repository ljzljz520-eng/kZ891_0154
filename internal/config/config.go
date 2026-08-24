package config

import "os"

type Config struct {
	Address        string
	DatabasePath   string
	AdminToken     string
	FixedToday     string
	SupportHotline string
	ServiceName    string
}

func Default() Config {
	return Config{Address: ":8080", DatabasePath: "warranty.db", AdminToken: "admin-fixed-token", FixedToday: "2025-06-01", SupportHotline: "400-800-1234", ServiceName: "家电安心延保服务"}
}

func Load() Config {
	c := Default()
	if v := os.Getenv("WARRANTY_ADDR"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("WARRANTY_DB"); v != "" {
		c.DatabasePath = v
	}
	if v := os.Getenv("WARRANTY_ADMIN_TOKEN"); v != "" {
		c.AdminToken = v
	}
	if v := os.Getenv("WARRANTY_TODAY"); v != "" {
		c.FixedToday = v
	}
	if v := os.Getenv("WARRANTY_HOTLINE"); v != "" {
		c.SupportHotline = v
	}
	if v := os.Getenv("WARRANTY_SERVICE_NAME"); v != "" {
		c.ServiceName = v
	}
	return c
}

func (c Config) Validate() error {
	if c.Address == "" {
		return ErrEmptyAddress
	}
	if c.DatabasePath == "" {
		return ErrEmptyDatabase
	}
	if c.AdminToken == "" {
		return ErrEmptyToken
	}
	if len(c.FixedToday) != 10 {
		return ErrInvalidDate
	}
	return nil
}

func (c Config) IsMemoryDatabase() bool {
	return c.DatabasePath == ":memory:" || c.DatabasePath == "file::memory:?cache=shared"
}
func (c Config) HotlineMessage() string { return c.ServiceName + "客服热线：" + c.SupportHotline }
