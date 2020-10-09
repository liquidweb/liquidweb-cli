package version

import (
	"fmt"
	"strings"

	"github.com/liquidweb/liquidweb-cli/utils"
)

// default build-time variables; these are overridden via ldflags
var (
	Version   = "unknown-version"
	GitCommit = "unknown-commit"
	BuildTime = "unknown-buildtime"
)

func Show() {
	fmt.Printf("LiquidWeb CLI Build Details\n\n")
	fmt.Printf("  Build Time: %s\n", BuildTime)
	fmt.Printf("  Version: %s\n", Version)
	fmt.Printf("  Git commit: %s\n\n", GitCommit)
}

func ShowLatest() error {
	runningLatest, running, latest, err := RunningLatestTag()
	if err != nil {
		return err
	}

	if runningLatest {
		utils.PrintGreen("Running latest tagged release\n")
	} else {
		if strings.Contains(running, "-dirty") {
			utils.PrintTeal("Running unofficial build [%s] latest [%s]\n", running, latest)
		} else {
			utils.PrintYellow("Update is available! running [%s] latest [%s]\n", running, latest)
		}
	}

	return nil
}
