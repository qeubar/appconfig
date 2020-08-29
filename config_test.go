package usrconfig

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanUpTestConfig(t *testing.T, app string) {
	fp, err := configFilePath(app)
	require.NoError(t, err)
	os.RemoveAll(fp)
}

func TestLoad(t *testing.T) {
	t.Run("handles no config file", func(t *testing.T) {
		type conf struct {
			User  string `xml:"user"`
			Email string `xml:"email"`
		}

		var c conf
		err := Load(&c, t.Name())
		assert.NoError(t, err, "did not ignore non-existent config")
		assert.Empty(t, c)
	})
	t.Run("returns unsupported struct", func(t *testing.T) {
		type conf struct {
			User  string `db:"user"`
			Email string `db:"email"`
		}

		var c conf
		err := Load(&c, t.Name())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported")
		assert.Empty(t, c)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("returns marshal err on bad struct", func(t *testing.T) {
		type conf struct {
			User  string      `xml:"user"`
			Email interface{} `xml:",comment"`
		}
		c := conf{
			User:  "bardia",
			Email: nil,
		}
		err := Update(c, t.Name())
		assert.Error(t, err)
	})
	t.Run("returns unsupported struct", func(t *testing.T) {
		type conf struct {
			User  string `db:"user"`
			Email string `db:"email"`
		}
		c := conf{
			User:  "bardia",
			Email: "bardia@keyoumarsi.com",
		}

		err := Update(c, t.Name())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported")
	})
}

func TestLoadAndUpdate(t *testing.T) {
	defer cleanUpTestConfig(t, t.Name())

	type conf struct {
		User  string `json:"user"`
		Email string `json:"email"`
	}

	c := conf{
		User:  "bardia",
		Email: "bardia@keyoumarsi.com",
	}

	err := Update(c, t.Name())
	require.NoError(t, err)

	var savedConf conf
	err = Load(&savedConf, t.Name())
	require.NoError(t, err)

	assert.Equal(t, c, savedConf)
}

func TestConfigFilePath(t *testing.T) {
	osCfgDir, err := os.UserConfigDir()
	require.NoError(t, err)

	fp, err := configFilePath("my-app")
	defer os.RemoveAll(filepath.Join(osCfgDir, "my-app"))

	require.NoError(t, err)
	assert.Equal(t, fp, filepath.Join(osCfgDir, "my-app", "config"))

	// Ensure that application can write to file
	err = ioutil.WriteFile(fp, nil, os.ModePerm)
	assert.NoError(t, err)
}

func TestConfigFileType(t *testing.T) {
	t.Run("detects json", func(t *testing.T) {
		type J struct {
			A string `json:"a"`
		}
		cft, err := configFileType(J{})
		require.NoError(t, err)
		assert.Equal(t, cftJSON, cft)
	})

	t.Run("detects yaml", func(t *testing.T) {
		type Y struct {
			A string `yaml:"a"`
		}
		cft, err := configFileType(Y{})
		require.NoError(t, err)
		assert.Equal(t, cftYAML, cft)
	})

	t.Run("detects xml", func(t *testing.T) {
		type X struct {
			A string `xml:"a"`
		}
		cft, err := configFileType(X{})
		require.NoError(t, err)
		assert.Equal(t, cftXML, cft)
	})

	t.Run("detects *struct", func(t *testing.T) {
		type J struct {
			A string `json:"a"`
		}
		cft, err := configFileType(&J{})
		require.NoError(t, err)
		assert.Equal(t, cftJSON, cft)
	})

	t.Run("errors on non-struct ", func(t *testing.T) {
		_, err := configFileType("qeubar")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a struct")
	})

	t.Run("errors on no-attribute struct", func(t *testing.T) {
		type N struct{}
		_, err := configFileType(N{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one field")
	})

	t.Run("errors on unsupported", func(t *testing.T) {
		type U struct {
			A string `unsupported:"bluh"`
		}
		_, err := configFileType(U{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported")
	})
}
