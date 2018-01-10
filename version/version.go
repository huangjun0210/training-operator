// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

import (
	"os"
	"runtime"

	log "github.com/golang/glog"
)

var (
	Version = "0.3.0+git"
	GitSHA  = "Not provided."
)

// PrintVersion print version info
func PrintVersion(shouldExit bool) {
	log.Infof("tf_operator Version: %v", Version)
	log.Infof("Git SHA: %s", GitSHA)
	log.Infof("Go Version: %s", runtime.Version())
	log.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	if shouldExit {
		os.Exit(0)
	}
}
