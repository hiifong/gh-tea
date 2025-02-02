package config

import (
	"fmt"

	pkgErr "github.com/hiifong/gh-tea/pkg/errors"
	"github.com/spf13/viper"
)

type TeaName string

func (t TeaName) String() string {
	return string(t)
}

type TeaItem struct {
	Name string
	Host string `yaml:"host"`
}

func (t TeaItem) String() string {
	return fmt.Sprintf(`
	Name: %s
	Host: %s
`, t.Name, t.Host)
}

var Default TeaName = "default"

type Tea map[TeaName]TeaItem

func (t Tea) String() string {
	var s string
	for k, v := range t {
		s += fmt.Sprintf("\n%s: %s", k, v)
	}
	return s
}

type Config struct {
	Tea Tea `yaml:"tea"`
}

func WriteConfig(Tea Tea) {
	viper.Set("tea", Tea)
	err := viper.WriteConfig()
	pkgErr.Check(err)
}
