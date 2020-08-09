package appconfig

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
