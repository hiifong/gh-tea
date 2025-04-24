package config

import (
	"fmt"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	pkgErr "github.com/hiifong/gh-tea/pkg/errors"
)

type (
	TeaName string

	Config struct {
		Tea Tea `yaml:"tea"`
	}

	Tea map[TeaName]TeaItem

	TeaItem struct {
		Name  string        `yaml:"name"`
		Host  string        `yaml:"host,omitempty"`
		Token *oauth2.Token `yaml:"token"`
	}
)

var (
	Default TeaName = "default"

	GiteaOAuth2 = oauth2.Config{
		ClientID: "d57cb8c4-630c-4168-8324-ec79935e18d4",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "%s/login/oauth/authorize",
			TokenURL: "%s/login/oauth/access_token",
		},
	}
)

func (t TeaName) String() string {
	return string(t)
}

func (t TeaItem) String() string {
	return fmt.Sprintf(`
	Name: %s
	Host: %s
	Token: %+v
`, t.Name, t.Host, t.Token)
}

func (t Tea) String() string {
	var s string
	for k, v := range t {
		s += fmt.Sprintf("\n%s: %s", k, v)
	}
	return s
}

func WriteConfig(Tea Tea) {
	viper.Set("tea", Tea)
	err := viper.WriteConfig()
	pkgErr.Check(err)
}
