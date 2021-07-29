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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var s3AccessKey, s3PrivateKey,
	s3Endpoint, s3Bucket,
	s3Region string
var s3ForcePathStyle bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "s3cli",
	Short: "s3cli to manipulate files on s3",
	Long:  `s3cli to manipulate files on s3`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.s3cli.yaml)")

	rootCmd.PersistentFlags().StringVar(&s3AccessKey, "accesskey", "none", "S3 Access Key")
	rootCmd.PersistentFlags().StringVar(&s3PrivateKey, "privatekey", "none", "S3 Private Key")
	rootCmd.PersistentFlags().StringVar(&s3Endpoint, "endpoint", "none", "S3 endpoint URL")
	rootCmd.PersistentFlags().StringVar(&s3Bucket, "bucket", "", "S3 Bucket")
	rootCmd.PersistentFlags().StringVar(&s3Region, "region", "us-east-1", "S3 Region")
	rootCmd.PersistentFlags().BoolVar(&s3ForcePathStyle, "forcepath", false, "If true, then will use s3 urls http://<endpoint>/<bucket>, otherwise http://<bucket>.<endpoint>/")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".s3cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".s3cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
