package config

import (
	"fmt"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	pkgErr "github.com/hiifong/gh-tea/pkg/errors"
)

type (
	Config struct {
		Tea Tea `yaml:"tea"`
	}
	Username string
	UserItem *oauth2.Token
	User     map[Username]UserItem

	TeaName string
	TeaItem struct {
		Name string
		Host string `yaml:"host"`
		User User   `yaml:"user"`
	}
	Tea map[TeaName]TeaItem
)

var (
	Default TeaName = "default"

	GiteaOAuth2 = oauth2.Config{
		ClientID: "9c0c7633-de85-4b51-938b-9fabc8cb7099",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://gitea.com/login/oauth/authorize",
			TokenURL: "https://gitea.com/login/oauth/access_token",
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
`, t.Name, t.Host)
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
