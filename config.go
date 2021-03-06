// Package usrconfig provides a minimal interface to store and update user specific
// application config on the supported platform. It takes away guess work and indecisiveness by
// using the default user config directory based on the running platform.
//
// It supports structs of `json`, `yaml` and `xml` format.
//
// It's as simple as:
//  type MyConfig struct {
//		Name  string `yaml:"user_name"`
//		Email string `yaml:"user_email"`
//  }
//
//  var conf MyConfig
//  usrconfig.Load(&conf, "my-app")
//
//  conf.Email = "updated@email.com"
//  usrconfig.Update(conf, "my-app")
package usrconfig

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-yaml/yaml"
)

const (
	cftJSON = "json"
	cftYAML = "yaml"
	cftXML  = "xml"
)

// Load will load an existing "config" for "app" saved under the
// user's config directory. If no config file exists for user,
// "config" remains empty and no error is returned.
func Load(config interface{}, app string) error {
	cft, err := configFileType(config)
	if err != nil {
		return err
	}

	configFilePath, err := configFilePath(app)
	if err != nil {
		return err
	}

	configBody, err := ioutil.ReadFile(configFilePath)
	switch {
	case err == nil:
		// NOOP
		break
	case os.IsNotExist(err):
		return nil // config doesn't exist, treat it as empty
	default:
		return err
	}

	switch cft {
	case cftJSON:
		err = json.Unmarshal(configBody, config)
	case cftYAML:
		err = yaml.Unmarshal(configBody, config)
	case cftXML:
		err = xml.Unmarshal(configBody, config)
	}

	return err
}

// Update encodes the provided "config" and saves it to the "app" config file.
// The file is written to the default user config directory on the platform.
func Update(config interface{}, app string) error {
	cft, err := configFileType(config)
	if err != nil {
		return err
	}

	configFilePath, err := configFilePath(app)
	if err != nil {
		return err
	}

	var cfgBody []byte
	switch cft {
	case cftJSON:
		cfgBody, err = json.MarshalIndent(config, "", "  ")
	case cftYAML:
		cfgBody, err = yaml.Marshal(config)
	case cftXML:
		cfgBody, err = xml.MarshalIndent(config, "", "  ")
	}
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFilePath, cfgBody, os.ModePerm)
}

func configFilePath(app string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(configDir, app), os.ModePerm)
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, app, "config"), nil
}

func configFileType(config interface{}) (string, error) {
	t := reflect.TypeOf(config)
	if t.Kind() == reflect.Ptr {
		// Extract the underlying object
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", errors.New("config must be a struct")
	}

	if t.NumField() < 1 {
		return "", errors.New("config must have at least one field")
	}

	for _, cft := range []string{cftJSON, cftYAML, cftXML} {
		if _, ok := t.Field(0).Tag.Lookup(cft); ok {
			return cft, nil
		}
	}

	return "", errors.New("unsupported config struct tag")
}
