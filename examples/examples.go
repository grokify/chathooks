package examples

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/grokify/simplego/os/osutil"
)

const (
	HandlersDir = "github.com/grokify/chathooks/docs/handlers"
	Examples    = "aha,appsignal,apteligent,circleci,codeship,confluence,datadog,deskdotcom,enchant,gosquared,heroku,librato,magnumci,marketo,opsgenie,papertrail,pingdom,raygun,runscope,semaphore,statuspage,travisci,userlike,victorops"
)

func AbsDirGopath(dir string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", dir)
}

func DocsHandlersDirInfo() ([]string, []string, error) {
	handlersDir := AbsDirGopath(HandlersDir)
	fmt.Println(handlersDir)

	dirs := []string{}
	exampleFiles := []string{}
	//sdirs, _, err := ioutilmore.ReadDirSplit(handlersDir, true)
	sdirs, err := osutil.ReadDirMore(handlersDir, nil, true, false, false)

	if err != nil {
		return dirs, exampleFiles, err
	}

	for _, sdir := range sdirs {
		fmt.Printf("SDIR: %v\n", sdir.Name())
		absSubDir := filepath.Join(handlersDir, sdir.Name())
		exEntries, err := osutil.ReadDirMore(absSubDir,
			regexp.MustCompile(`^event-example_.+\.(json|txt)$`), false, true, false)
		if err != nil {
			return dirs, exampleFiles, err
		}
		if len(exEntries) > 0 {
			dirs = append(dirs, sdir.Name())
			for _, f := range exEntries {
				fmt.Printf("FILE: %v\n", f.Name())
				exFilepath := filepath.Join(absSubDir, f.Name())
				exampleFiles = append(exampleFiles, exFilepath)
			}
		}
	}
	return dirs, exampleFiles, nil
}
