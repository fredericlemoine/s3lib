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

var deletePath string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete files from s3",
	Long:  `delete files from s3`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if deletePath == "none" {
			err = fmt.Errorf("A path must be specified to delete a file")
			return
		}

		var cli = s3lib.NewS3Cli(s3AccessKey, s3PrivateKey, s3Endpoint, s3Region, s3ForcePathStyle)

		if err = cli.Delete(s3Bucket, deletePath); err != nil {
			return
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deletePath, "s3path", "none", "S3 Path to the file to delete")
}
