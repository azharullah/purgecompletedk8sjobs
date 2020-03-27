/*
Copyright © 2020 Azharullah mdazharullah.shariff@nutanix.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"

	"github.com/azharullah/purge-completed-k8s-jobs/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	cmd.Execute()
}

func init() {

	// Set the log level, if not specified via env
	lvl, ok := os.LookupEnv("LOG_LEVEL")

	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}

	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}

	// set global log level
	logrus.SetLevel(ll)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		// CallerPrettyfier: func(f *runtime.Frame) (string, string) {
		// 	filename := path.Base(f.File)
		// 	log.Printf("%#v", f)
		// 	return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		// },
	})
}
