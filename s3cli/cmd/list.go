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
	"time"

	"github.com/fredericlemoine/s3lib"
	"github.com/spf13/cobra"
)

var listPrefix string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list files from s3",
	Long:  `list files from s3`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var cli = s3lib.NewS3Cli(s3AccessKey, s3PrivateKey, s3Endpoint, s3Region, s3ForcePathStyle)
		var files []s3lib.S3File

		if files, err = cli.List(s3Bucket, listPrefix); err != nil {
			return
		}

		for _, f := range files {
			fmt.Print(f.LastModified.Format(time.RFC822))
			fmt.Printf("\t%d", f.Size)
			fmt.Printf("\t%s\n", f.Key)
			//fmt.Printf("\t%s\n", f.StorageClass)
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listPrefix, "s3prefix", "", "Prefix of the files to list")
}
