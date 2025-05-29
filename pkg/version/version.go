// The MIT License (MIT)
//
// Copyright © 2025 Yusheng Guo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package version

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

var (
	// version is a constant representing the version tag that
	// generated this build. It should be set during build via -ldflags.
	version string
	// buildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	// It should be set during build via -ldflags.
	buildDate string
	// gitbranch is a constant representing git branch for this build.
	// It should be set during build via -ldflags.
	gitbranch string
	// gitbranch is a constant representing git sha1 for this build.
	// It should be set during build via -ldflags.
	gitsha1 string
)

// Info holds the information related to kle app version.
type Info struct {
	Major      string `json:"major"`
	Minor      string `json:"minor"`
	GitVersion string `json:"gitVersion"`
	GitBranch  string `json:"gitBranch"`
	GitSha1    string `json:"gitSha1"`
	BuildDate  string `json:"buildDate"`
	GoVersion  string `json:"goVersion"`
	Compiler   string `json:"compiler"`
	Platform   string `json:"platform"`
}

func Get() Info {
	majorVersion, minorVersion := splitVersion(version)
	return Info{
		Major:      majorVersion,
		Minor:      minorVersion,
		GitVersion: version,
		GitBranch:  gitbranch,
		GitSha1:    gitsha1,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// splitVersion splits the git version to generate major and minor versions needed.
func splitVersion(version string) (string, string) {
	if version == "" {
		return "", ""
	}

	// Version from an automated container build environment for a tag. For example v20200521-v0.18.0.
	m1, _ := regexp.MatchString(`^v\d{8}-v\d+\.\d+\.\d+$`, version)

	// Version from an automated container build environment(not a tag) or a local build. For example v20201009-v0.18.0-46-g939c1c0.
	m2, _ := regexp.MatchString(`^v\d{8}-v\d+\.\d+\.\d+-\w+-\w+$`, version)

	// Version tagged by helm chart releaser action
	helm, _ := regexp.MatchString(`^v\d{8}-descheduler-helm-chart-\d+\.\d+\.\d+$`, version)
	// Dirty version where helm chart is the last known tag
	helm2, _ := regexp.MatchString(`^v\d{8}-descheduler-helm-chart-\d+\.\d+\.\d+-\w+-\w+$`, version)

	if m1 || m2 {
		semVer := strings.Split(version, "-")[1]
		return strings.Trim(strings.Split(semVer, ".")[0], "v"), strings.Split(semVer, ".")[1] + "." + strings.Split(semVer, ".")[2]
	}

	if helm || helm2 {
		semVer := strings.Split(version, "-")[4]
		return strings.Split(semVer, ".")[0], strings.Split(semVer, ".")[1] + "." + strings.Split(semVer, ".")[2]
	}

	// Something went wrong
	return "", ""
}
