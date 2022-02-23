/*
Copyright 2016 The Kubernetes Authors.

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
	"domain-assets/pkg/cmdutils"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Flags: cmdutils.MergeFlags(
			azureFlags(),
			awsFlags(),
		),
		Action: func(c *cli.Context) error {
			go func() {
				initApp(c)
				awsAssets := runAWS(c)
				azureAssets := runAzure(c)
				DNSAssets := append(awsAssets, azureAssets...)
				manageAssets(DNSAssets)
				time.Sleep(time.Duration(c.Uint("sync-loop")) * time.Second)

			}()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
