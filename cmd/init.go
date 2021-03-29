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
	"errors"
	"os"
	"path"

	"github.com/devrel-blox/drb/config"
	"github.com/spf13/cobra"
)

var (
	base        string
	source      string
	destination string
	templates   string
	skipConfig  bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create folders and configuration to maintain content with the drb toolset",
	Long: `Create a group of folders to store your content. 

If provided, the folders will be created under the "base" directory. 
If "base" is set to an empty string, the source, destination, and template
folders will be created in the root of the current directory.

The "source" directory will store your un-processed content, 
typically Markdown files.

The "destination" directory is where the drb tools will put 
content after it has been validated and processed.

The "template" directory is where you can store templates for
each content type with pre-filled values.
`,
	Run: func(cmd *cobra.Command, args []string) {

		root, err := os.Getwd()
		if err != nil {
			cmd.PrintErr("unable to get current directory")
			return
		}
		err = createDirectories(root)
		if err != nil {
			cmd.PrintErr(err.Error())
			return
		}

		if !skipConfig {
			err = writeConfigFile()
			if err != nil {
				cmd.PrintErr(err.Error())
				return
			}
		}

	},
}

func writeConfigFile() error {
	cfg := config.BloxConfig{
		Base:        base,
		Source:      source,
		Templates:   templates,
		Destination: destination,
	}
	f, err := os.Create("blox.yaml")
	if err != nil {
		return err
	}
	defer f.Close()
	err = cfg.Write(f)
	return err
}

func createDirectories(root string) error {
	err := os.MkdirAll(sourceDir(root), 0755)
	if err != nil {
		return errors.New("error creating source directory")
	}
	err = os.MkdirAll(destinationDir(root), 0755)
	if err != nil {
		return errors.New("error creating destination directory")
	}
	err = os.MkdirAll(templateDir(root), 0755)
	if err != nil {
		return errors.New("error creating template directory")
	}

	return nil
}

func sourceDir(root string) string {
	return path.Join(root, base, source)
}
func destinationDir(root string) string {
	return path.Join(root, base, destination)
}
func templateDir(root string) string {
	return path.Join(root, base, templates)
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&base, "base", "b", "content", "base directory for pre- and post- processed content")
	initCmd.Flags().StringVarP(&source, "source", "s", "source", "where pre-processed content will be stored (source markdown)")
	initCmd.Flags().StringVarP(&destination, "destination", "d", "out", "where post-processed content will be stored (output json)")
	initCmd.Flags().StringVarP(&templates, "template", "t", "templates", "where content templates will be stored")
	initCmd.Flags().BoolVarP(&skipConfig, "skip", "c", false, "don't write a configuration file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
