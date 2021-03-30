/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"path"

	"github.com/devrel-blox/drb/models"
	"github.com/spf13/cobra"
)

var (
	model string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new content file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model, err := models.GetModel(model)
		cobra.CheckErr(err)

		fmt.Printf("Creating new %s in %s\n", model.Name, model.SourceContentPath())

		slug := args[0]

		cobra.CheckErr(model.New(slug))
		fmt.Printf("Your new content file is ready at %s\n", path.Join(model.SourceFilePath(slug)))
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&model, "type", "t", "profile", "type of content to create")
	cobra.CheckErr(newCmd.MarkFlagRequired("type"))
	newCmd.SetUsageTemplate("drb new --type [type name] [slug]")
}
