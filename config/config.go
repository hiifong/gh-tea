package config

type TeaName string

func (t TeaName) String() string {
	return string(t)
}

type Tea struct {
	URL string `yaml:"url"`
}

type Config struct {
	Tea map[TeaName]Tea `yaml:"tea"`
}
