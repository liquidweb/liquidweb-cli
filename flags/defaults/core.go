package defaults

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"github.com/liquidweb/liquidweb-cli/config"
	"github.com/liquidweb/liquidweb-cli/utils"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var (
	nagged map[string]bool
	nags   bool
	tipped bool
)

func init() {
	nagged = map[string]bool{}

	home, err := homedir.Dir()
	if err != nil {
		utils.PrintYellow("failed fetching homedir: %s\n", err)
		return
	}
	viper.SetDefault("liquidweb.flags.defaults.file", fmt.Sprintf("%s/.liquidweb-cli-flag-defaults.yaml", home))
}

func NagsOff() (err error) {
	err = toggleNags(false)

	return
}

func NagsOn() (err error) {
	err = toggleNags(true)

	return
}

func GetPermitted() (permitted []string) {
	permitted = make([]string, 0, len(permittedFlags))
	for flag, opts := range permittedFlags {
		if enabled, ok := opts["enabled"].(bool); !enabled || !ok {
			continue
		}
		permitted = append(permitted, flag)
	}
	sort.Strings(permitted)

	return
}

func GetOrNag(flag string) (value interface{}) {
	var err error
	value, err = Get(flag)
	if err != nil {
		if !nagged[flag] {
			if errors.Is(err, ErrorNotFound) {
				if nags {
					fmt.Printf("No default value for flag [%s] set. See 'help default-flags set' for details.\n", flag)
					if !tipped {
						utils.PrintTeal("TIP: You can silence undefined default flag notices with 'default-flags nags-off'\n")
						tipped = true
					}
				}
			} else {
				utils.PrintYellow("WARNING: Unexpected error when fetching value for default flag [%s]: %s\n", flag, err)
			}
			nagged[flag] = true
		}
	}
	return
}

func Get(flag string) (value interface{}, err error) {
	if _, err = permittedFlagOrError(flag); err != nil {
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
	var validator string
	validator, err = permittedFlagOrError(flag)
	if err != nil {
		return
	}

	if err = validateFlagValue(validator, value); err != nil {
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
	vp.Set(contextFlagKey(), flags)

	err = writeViperConfig(vp)

	return
}

func Delete(flag string) (err error) {
	if _, err = permittedFlagOrError(flag); err != nil {
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
	vp.Set(contextFlagKey(), flags)
	err = writeViperConfig(vp)

	return
}

func permittedFlagOrError(flag string) (validator string, err error) {
	if flag == "" {
		err = ErrorInvalidFlagName
		return
	}

	opts, exists := permittedFlags[flag]
	if !exists {
		err = fmt.Errorf("%s %w", flag, ErrorForbiddenFlag)
		return
	}

	if v, ok := opts["enabled"].(bool); !ok || !v {
		err = fmt.Errorf("%s %w", flag, ErrorForbiddenFlag)
	}

	validator = cast.ToString(opts["type"])

	return
}

func validateFlagValue(validator string, value interface{}) (err error) {
	var validateFields map[interface{}]interface{}

	if strings.HasSuffix(validator, "Int64") {
		validateFields = map[interface{}]interface{}{
			cast.ToInt64(value): "PositiveInt64",
		}
	} else if strings.HasSuffix(validator, "Int") {
		validateFields = map[interface{}]interface{}{
			cast.ToInt(value): "PositiveInt64",
		}
	} else if strings.HasSuffix(validator, "String") {
		validateFields = map[interface{}]interface{}{
			cast.ToString(value): "NonEmptyString",
		}
	}

	err = validate.Validate(validateFields)

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

	flags = vp.GetStringMap(contextFlagKey())

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
	vp.SetDefault(NagsKey, true)
	if err = vp.ReadInConfig(); err != nil {
		err = fmt.Errorf("%w: %s", ErrorUnreadable, err)
		return
	}

	nags = vp.GetBool(NagsKey)

	return
}

func getFlagsFile() (file string, err error) {
	file = viper.GetString(DefaultFlagsFileKey)
	if file == "" {
		err = ErrorFileKeyMissing
	}

	if _, err = os.Stat(file); os.IsNotExist(err) {
		err = nil
		f, ferr := os.Create(filepath.Clean(file))
		if ferr != nil {
			err = ferr
			return
		}
		err = f.Close()
	}

	return
}

func contextFlagKey() (k string) {
	k = fmt.Sprintf("%s.%s", DefFlagsKey, config.CurrentContext)

	return
}

func toggleNags(on bool) error {
	vp, err := getFlagsViper()
	if err != nil {
		return err
	}

	vp.Set(NagsKey, on)

	err = writeViperConfig(vp)

	return err
}

func writeViperConfig(vp *viper.Viper) (err error) {
	if err = vp.WriteConfig(); err != nil {
		err = fmt.Errorf("%w: %s", ErrorUnwritable, err)
	}

	return
}
