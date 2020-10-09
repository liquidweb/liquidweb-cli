package version

import (
	"context"

	"github.com/google/go-github/v32/github"
)

func GetLatestTag() (tag string, err error) {
	client := github.NewClient(nil)

	tags, _, err := client.Repositories.ListTags(context.Background(), "liquidweb", "liquidweb-cli", nil)
	if err != nil {
		return
	}

	latest := tags[0]
	tag = *latest.Name

	return
}

func RunningLatestTag() (runningLatest bool, running, latest string, err error) {
	latest, err = GetLatestTag()
	if err != nil {
		return
	}

	running = Version

	if running == latest {
		runningLatest = true
	}

	return
}
