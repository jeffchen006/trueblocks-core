// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package config

import (
	"log"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/usage"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/version"
	"github.com/spf13/viper"
)

var trueBlocksViper = viper.New()
var trueBlocksConfig ConfigFile

type versionGroup struct {
	Current string `toml:"current"`
}

type chainGroup struct {
	ChainId        string `toml:"chainId"`
	LocalExplorer  string `toml:"localExplorer"`
	RemoteExplorer string `toml:"remoteExplorer"`
	RpcProvider    string `toml:"rpcProvider"`
	ApiProvider    string `toml:"apiProvider"`
	IpfsGateway    string `toml:"ipfsGateway"`
	Symbol         string `toml:"symbol"`
}

type keyGroup struct {
	ApiKey string `toml:"apiKey"`
	Secret string `toml:"secret"`
	Jwt    string `toml:"jwt"`
}

type settingsGroup struct {
	CachePath      string `toml:"cachePath"`
	IndexPath      string `toml:"indexPath"`
	DefaultChain   string `toml:"defaultChain"`
	DefaultGateway string `toml:"defaultGateway"`
}

type ConfigFile struct {
	Version  versionGroup
	Settings settingsGroup
	Keys     map[string]keyGroup
	Chains   map[string]chainGroup
}

// init sets up default values for the given configuration
func init() {
	trueBlocksViper.SetConfigName("trueBlocks") // trueBlocks.toml (so we can find it)
	trueBlocksViper.SetDefault("Settings.CachePath", GetPathToRootConfig()+"cache/")
	trueBlocksViper.SetDefault("Settings.IndexPath", GetPathToRootConfig()+"unchained/")
	trueBlocksViper.SetDefault("Settings.DefaultChain", "mainnet")
	trueBlocksViper.SetDefault("Settings.DefaultGateway", "https://ipfs.unchainedindex.io/ipfs")
}

// GetRootConfig reads and the configuration located in trueBlocks.toml file. Note
// that this routine is local to the package
func GetRootConfig() *ConfigFile {
	if len(trueBlocksConfig.Settings.CachePath) == 0 {
		configPath := GetPathToRootConfig()
		MustReadConfig(trueBlocksViper, &trueBlocksConfig, configPath)

		user, _ := user.Current()

		cachePath := trueBlocksConfig.Settings.CachePath
		if len(cachePath) == 0 || cachePath == "<not_set>" {
			cachePath = path.Join(configPath, "cache") + "/"
		}
		cachePath = strings.Replace(cachePath, "$HOME", user.HomeDir, -1)
		cachePath = strings.Replace(cachePath, "~", user.HomeDir, -1)
		trueBlocksConfig.Settings.CachePath = cachePath

		indexPath := trueBlocksConfig.Settings.IndexPath
		if len(indexPath) == 0 || indexPath == "<not_set>" {
			indexPath = path.Join(configPath, "unchained") + "/"
		}
		indexPath = strings.Replace(indexPath, "$HOME", user.HomeDir, -1)
		indexPath = strings.Replace(indexPath, "~", user.HomeDir, -1)
		trueBlocksConfig.Settings.IndexPath = indexPath

		if len(trueBlocksConfig.Settings.DefaultChain) == 0 {
			trueBlocksConfig.Settings.DefaultChain = "mainnet"
		}

		// We establish only the top-level folders here. When we figure out
		// which chain we're on (not until the user tells us on the command line)
		// only then can we complete these paths. At this point these paths
		// only point to the top-levl of the cache or index. Also note that
		// these two calls do not return if they fail, so no need to handle errors
		defaultChains := []string{GetDefaultChain()}
		file.EstablishFolders(trueBlocksConfig.Settings.CachePath, defaultChains)
		file.EstablishFolders(trueBlocksConfig.Settings.IndexPath, defaultChains)
	}

	return &trueBlocksConfig
}

func IsAtLeastVersion(needle string) bool {
	var current, desired version.Version
	var err error
	if current, err = version.NewVersion(GetRootConfig().Version.Current); err != nil {
		return true
	}

	if desired, err = version.NewVersion(needle); err != nil {
		return true
	}

	return !current.IsEarlierThan(desired)
}

// GetPathToRootConfig returns the path where to find configuration files
func GetPathToRootConfig() string {
	configPath, err := PathFromXDG("XDG_CONFIG_HOME")
	if err != nil {
		log.Fatal(err)
	} else if len(configPath) > 0 {
		return configPath
	}

	// The migration code will have already checked for invalid operating systems (i.e. Windows)
	userOs := runtime.GOOS
	if len(os.Getenv("TEST_OS")) > 0 {
		userOs = os.Getenv("TEST_OS")
	}

	user, _ := user.Current()
	osPath := ".local/share/trueblocks"
	if userOs == "darwin" {
		osPath = "Library/Application Support/TrueBlocks"
	}

	return path.Join(user.HomeDir, osPath) + "/"
}

func GetDefaultChain() string {
	return GetRootConfig().Settings.DefaultChain
}

func PathFromXDG(envVar string) (string, error) {
	// If present, we require both an existing path and a fully qualified path
	xdg := os.Getenv(envVar)
	if len(xdg) == 0 {
		return "", nil // it's okay if it's empty
	}

	if xdg[0] != '/' {
		return "", usage.Usage("The {0} value ({1}), must be fully qualified.", envVar, xdg)
	}

	if _, err := os.Stat(xdg); err != nil {
		return "", usage.Usage("The {0} folder ({1}) must exist.", envVar, xdg)
	}

	return path.Join(xdg, "") + "/", nil
}

// MustReadConfig calls Viper's ReadInConfig and fills values in the
// given targetStruct. Any error will result in a call to logger.Fatal
func MustReadConfig(v *viper.Viper, targetStruct interface{}, path string) {
	v.AddConfigPath(path)
	v.SetEnvPrefix("TB")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		// We only require some files to be present
		if !ok {
			logger.Fatal(err)
		}
	}

	err = v.Unmarshal(targetStruct)
	if err != nil {
		logger.Fatal(err)
	}
}

func GetPinningKeys(chain string) (string, string, string) {
	keys := GetRootConfig().Keys
	a := keys["pinata"].ApiKey
	b := keys["pinata"].Secret
	c := keys["estuary"].ApiKey
	return a, b, c
}
