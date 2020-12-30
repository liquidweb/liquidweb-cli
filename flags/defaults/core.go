package defaults

import (
	"errors"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/liquidweb/liquidweb-cli/utils"
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		utils.PrintYellow("failed fetching homedir: %s\n", err)
		return
	}
	viper.SetDefault("liquidweb.flags.defaults.file", fmt.Sprintf("%s/.liquidweb-cli-flag-defaults.yaml", home))
}

//var flagDefaultsFile = viper.GetString("liquidweb.flags.defaults.file")

func GetOrNag(flag string) (value interface{}) {
	var err error
	value, err = Get(flag)
	if err != nil {
		if errors.Is(err, ErrorNotFound) {
			utils.PrintTeal("No default for flag [%s] set. See 'help default-flags set' for details.\n", flag)
		} else {
			utils.PrintYellow("Unexpected error when fetching value for default flag [%s]: %s\n", flag, err)
		}
	}
	return
}

func Get(flag string) (value interface{}, err error) {
	if err = permittedFlagOrError(flag); err != nil {
		return
	}

	var flags map[string]interface{}
	flags, err = getFlagsMap()
	if err != nil {
		return
	}

	if v, exists := flags[flag]; exists {
		value = v
		return
	}

	err = fmt.Errorf("%s %w", flag, ErrorNotFound)
	return
}

func GetAll() (all AllFlags, err error) {
	all, err = getFlagsMap()

	return
}

func Set(flag string, value interface{}) (err error) {
	if err = permittedFlagOrError(flag); err != nil {
		return
	}

	var (
		vp    *viper.Viper
		flags map[string]interface{}
	)
	vp, flags, err = getFlagsViperAndMap()
	if err != nil {
		return
	}

	flags[flag] = value
	vp.Set(DefFlagsKey, flags)

	if err = vp.WriteConfig(); err != nil {
		err = fmt.Errorf("%w: %s", ErrorUnwritable, err)
	}

	return
}

func Delete(flag string) (err error) {
	if err = permittedFlagOrError(flag); err != nil {
		return
	}

	var (
		vp    *viper.Viper
		flags map[string]interface{}
	)
	vp, flags, err = getFlagsViperAndMap()
	if err != nil {
		return
	}

	delete(flags, flag)
	vp.Set(DefFlagsKey, flags)
	err = vp.WriteConfig()

	return
}

func permittedFlagOrError(flag string) (err error) {
	if flag == "" {
		err = ErrorInvalidFlagName
		return
	}

	if _, exists := permittedFlags[flag]; !exists {
		err = fmt.Errorf("%s %w", flag, ErrorForbiddenFlag)
	}

	return
}

func getFlagsViperAndMap() (vp *viper.Viper, flags map[string]interface{}, err error) {
	vp, err = getFlagsViper()
	if err != nil {
		return
	}

	flags, err = getFlagsMap(vp)

	return
}

func getFlagsMap(vpL ...*viper.Viper) (flags map[string]interface{}, err error) {
	var vp *viper.Viper
	if len(vpL) == 0 {
		if vp, err = getFlagsViper(); err != nil {
			return
		}
	} else {
		vp = vpL[0]
	}

	flags = vp.GetStringMap(DefFlagsKey)

	return
}

func getFlagsViper() (vp *viper.Viper, err error) {
	var file string
	file, err = getFlagsFile()
	if err != nil {
		return
	}

	vp = viper.New()
	vp.SetConfigFile(file)
	if err = vp.ReadInConfig(); err != nil {
		err = fmt.Errorf("%w: %s", ErrorUnreadable, err)
		return
	}

	return
}

func getFlagsFile() (file string, err error) {
	file = viper.GetString(DefaultFlagsFileKey)
	if file == "" {
		err = ErrorFileKeyMissing
	}

	if _, err = os.Stat(file); os.IsNotExist(err) {
		err = nil
		f, ferr := os.Create(file)
		if ferr != nil {
			err = ferr
			return
		}
		err = f.Close()
	}

	return
}
