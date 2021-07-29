/*
Copyright © 2021 Institut Pasteur

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
package cmd

import (
	"fmt"

	"github.com/fredericlemoine/s3lib"
	"github.com/spf13/cobra"
)

var downloadFilename, downloadPath string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads a file from s3",
	Long:  `Downloads a file from s3`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if downloadFilename == "none" {
			err = fmt.Errorf("outfile must be specified")
			return
		}
		var cli = s3lib.NewS3Cli(s3AccessKey, s3PrivateKey, s3Endpoint, s3Region, s3ForcePathStyle)
		err = cli.Download(downloadFilename, s3Bucket, downloadPath)
		return
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVar(&downloadFilename, "outfile", "none", "Downloaded local file name")
	downloadCmd.Flags().StringVar(&downloadPath, "s3path", "none", "S3 path to the file to download")
}
