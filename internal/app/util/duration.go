package util

import "time"

// Duration — обёртка над time.Duration для парсинга из ENV и JSON.
type Duration struct {
	time.Duration
}

// UnmarshalText реализует encoding.TextUnmarshaler для envconfig.
func (d *Duration) UnmarshalText(text []byte) error {
	parsed, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}
