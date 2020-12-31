package examples

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/grokify/simplego/io/ioutilmore"
)

const (
	HandlersDir = "github.com/grokify/chathooks/docs/handlers"
	Examples    = "aha,appsignal,apteligent,circleci,codeship,confluence,datadog,deskdotcom,enchant,gosquared,heroku,librato,magnumci,marketo,opsgenie,papertrail,pingdom,raygun,runscope,semaphore,statuspage,travisci,userlike,victorops"
)

func AbsDirGopath(dir string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", dir)
}

func DocsHandlersDirInfo() ([]string, []string, error) {
	dirname := AbsDirGopath(HandlersDir)
	fmt.Println(dirname)

	dirs := []string{}
	exampleFiles := []string{}
	sdirs, _, err := ioutilmore.ReadDirSplit(dirname, true)

	if err != nil {
		return dirs, exampleFiles, err
	}

	for _, sdir := range sdirs {
		fmt.Printf("SDIR: %v\n", sdir.Name())
		absSubDir := filepath.Join(dirname, sdir.Name())
		files, _, err := ioutilmore.ReadDirRx(absSubDir,
			regexp.MustCompile(`^event-example_.+\.(json|txt)$`), true)
		if err != nil {
			return dirs, exampleFiles, err
		}
		if len(files) > 0 {
			dirs = append(dirs, sdir.Name())
			for _, f := range files {
				fmt.Printf("FILE: %v\n", f.Name())
				exFilepath := filepath.Join(absSubDir, f.Name())
				exampleFiles = append(exampleFiles, exFilepath)
			}
		}
	}
	return dirs, exampleFiles, nil
}
