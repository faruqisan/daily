package secret

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type (
	// Secret struct hold all secret config related
	Secret struct {
		GoogleOAuth GoogleOAuth `yaml:"google_oauth"`
	}

	// GoogleOAuth struct hold all google oauth secret
	GoogleOAuth struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectURL  string `yaml:"redirect_url"`
	}
)

var sec *Secret

// Get function return ready to use secret obj
func Get() (*Secret, error) {
	if sec != nil {
		return sec, nil
	}

	file, err := ioutil.ReadFile("files/secret.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &sec)
	if err != nil {
		return nil, err
	}

	return sec, nil
}
