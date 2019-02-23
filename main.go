package main

import (
	"net/http"

	"github.com/Luzifer/go_helpers/str"
	"github.com/Luzifer/nginx-sso/plugins"

	yaml "gopkg.in/yaml.v2"
)

func Register(ra plugins.RegisterAuthenticatorFunc, rm plugins.RegisterMFAProviderFunc) {
	ra(&authSupercookie{})
}

type authSupercookie struct {
	CookieName string              `yaml:"cookie_name"`
	Tokens     map[string]string   `yaml:"tokens"`
	Groups     map[string][]string `yaml:"groups"`
}

// AuthenticatorID needs to return an unique string to identify
// this special authenticator
func (a authSupercookie) AuthenticatorID() string { return "supercookie" }

// Configure loads the configuration for the Authenticator from the
// global config.yaml file which is passed as a byte-slice.
// If no configuration for the Authenticator is supplied the function
// needs to return the plugins.ErrProviderUnconfigured
func (a *authSupercookie) Configure(yamlSource []byte) error {
	envelope := struct {
		Providers struct {
			SuperCookie *authSupercookie `yaml:"supercookie"`
		} `yaml:"providers"`
	}{}

	if err := yaml.Unmarshal(yamlSource, &envelope); err != nil {
		return err
	}

	if envelope.Providers.SuperCookie == nil {
		return plugins.ErrProviderUnconfigured
	}

	a.CookieName = envelope.Providers.SuperCookie.CookieName
	a.Tokens = envelope.Providers.SuperCookie.Tokens
	a.Groups = envelope.Providers.SuperCookie.Groups

	return nil
}

// DetectUser is used to detect a user without a login form from
// a cookie, header or other methods
// If no user was detected the plugins.ErrNoValidUserFound needs to be
// returned
func (a authSupercookie) DetectUser(res http.ResponseWriter, r *http.Request) (string, []string, error) {
	var (
		user, token string
		userFound   bool
	)

	for _, c := range r.Cookies() {
		if c.Name != a.CookieName {
			continue
		}

		for user, token = range a.Tokens {
			if token == c.Value {
				userFound = true
				break
			}
		}

		if userFound {
			break
		}
	}

	if !userFound {
		return "", nil, plugins.ErrNoValidUserFound
	}

	groups := []string{}
	for group, users := range a.Groups {
		if str.StringInSlice(user, users) {
			groups = append(groups, group)
		}
	}

	return user, groups, nil
}

// Login is called when the user submits the login form and needs
// to authenticate the user or throw an error. If the user has
// successfully logged in the persistent cookie should be written
// in order to use DetectUser for the next login.
// If the user did not login correctly the plugins.ErrNoValidUserFound
// needs to be returned
func (a authSupercookie) Login(res http.ResponseWriter, r *http.Request) (string, []plugins.MFAConfig, error) {
	return "", nil, plugins.ErrNoValidUserFound
}

// LoginFields needs to return the fields required for this login
// method. If no login using this method is possible the function
// needs to return nil.
func (a authSupercookie) LoginFields() []plugins.LoginField { return nil }

// Logout is called when the user visits the logout endpoint and
// needs to destroy any persistent stored cookies
func (a authSupercookie) Logout(res http.ResponseWriter, r *http.Request) error { return nil }

// SupportsMFA returns the MFA detection capabilities of the login
// provider. If the provider can provide mfaConfig objects from its
// configuration return true. If this is true the login interface
// will display an additional field for this provider for the user
// to fill in their MFA token.
func (a authSupercookie) SupportsMFA() bool { return false }
