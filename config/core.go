package config

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/utils"
)

var (
	CurrentContext string
	ConfigFileArg  string
	UseContextArg  string
)

func InitConfig() (vp *viper.Viper, err error) {
	vp = viper.New()

	if ConfigFileArg != "" {
		// Use config file from the flag.
		vp.SetConfigFile(ConfigFileArg)
	} else {
		// Search config in home directory with name ".liquidweb-cli" (without extension).
		var home string
		home, err = homedir.Dir()
		if err != nil {
			return
		}
		vp.AddConfigPath(home)
		vp.SetConfigName(".liquidweb-cli")
	}

	vp.AutomaticEnv()
	if err = vp.ReadInConfig(); err != nil {
		if _, notFound := err.(viper.ConfigFileNotFoundError); notFound {
			err = nil
			return
		}
		utils.PrintYellow("error reading config: %s\n", err)
		return
	}

	if UseContextArg != "" {
		if err = instance.ValidateContext(UseContextArg, vp); err != nil {
			err = fmt.Errorf("error using auth context: %s\n", err)
			return
		}
		vp.Set("liquidweb.api.current_context", UseContextArg)
	}

	CurrentContext = vp.GetString("liquidweb.api.current_context")

	return
}
