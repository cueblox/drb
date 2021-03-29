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

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		/*path, err := os.Getwd()

		if err != nil {
			log.Println(err)
			return
		}
		*/

		var runtime cue.Runtime

		// The entrypoints are the same as the files you'd specify at the command line
		//entrypoints := []string{fmt.Sprintf("%s/.cache/cue/root.cue", path)}

		entrypoints := []string{"github.com/devrel-blox/schema/profile"}
		// Load Cue files into Cue build.Instances slice
		// the second arg is a configuration object, we'll see this later

		thingies := load.Instances(entrypoints, nil)

		for _, bi := range thingies {
			// check for errors on the instance
			// these are typically parsing errors
			if bi.Err != nil {
				fmt.Println("Error during load:", bi.Err)
				continue
			}

			// Use cue.Runtime to build.Instance to cue.INstance
			I, err := runtime.Build(bi)
			if err != nil {
				fmt.Println("Error during build:", bi.Err)
				continue
			}

			// get the root value and print it
			value := I.Value()
			fmt.Println("root value:", value)

			// Validate the value
			err = value.Validate()
			if err != nil {
				fmt.Println("Error during validate:", err)
				continue
			}
		}

		// spew.Dump(buildConfig)

		// instances := load.Instances(nil, buildConfig)

		// c := cue.Build(instances)

		// for _, j := range c {
		// 	fmt.Println("Found: ", j.Eval())
		// }

		fmt.Println("validate called")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
