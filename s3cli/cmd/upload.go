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
	"github.com/fredericlemoine/s3lib"
	"github.com/spf13/cobra"
)

var uploadFilename string
var uploadPath string

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to s3",
	Long:  `Upload a file to s3`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var cli = s3lib.NewS3Cli(s3AccessKey, s3PrivateKey, s3Endpoint, s3Region, s3ForcePathStyle)
		err = cli.Upload(uploadFilename, s3Bucket, uploadPath)
		return
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.
	uploadCmd.Flags().StringVar(&uploadFilename, "infile", "none", "File to upload")
	uploadCmd.Flags().StringVar(&uploadPath, "s3path", "none", "Path to upload the file on S3")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
