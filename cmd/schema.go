package cmd

// import (
// 	_ "embed"
// 	"fmt"

// 	"cuelang.org/go/cue"
// 	"github.com/spf13/cobra"
// )

// //go:embed schema_validation.cue
// var schemaValidation string

// //go:embed schema.cue
// var downloadedCue string

// var schemaCmd = &cobra.Command{
// 	Use:   "schema",
// 	Short: "",
// 	Long:  ``,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// This function is assuming we've downloaded a schema
// 		// to the schemas directory. It will now attempt to load
// 		// it, validate it, and extract the models
// 		var cueRuntime cue.Runtime
// 		remoteSchema, err := cueRuntime.Compile("schema", downloadedCue)
// 		cobra.CheckErr(err)

// 		remoteSchemaValue := remoteSchema.Value()

// 		// This will return a struct with whatever we've got defined
// 		// as part of our schema metadata, in schema_validation.cue's #Schema
// 		schemaStruct, err := validateSchema(remoteSchemaValue)
// 		if err != nil {
// 			fmt.Println("Remote schema is not valid")
// 			cobra.CheckErr(err)
// 		}

// 		// Confirm namespace and name are available
// 		fmt.Println(schemaStruct)

// 		// root := cue.Path{}

// 		// // Is an acceptable / passed schema definition,
// 		// // can we find the values?
// 		fields, err := remoteSchemaValue.Fields(cue.All())
// 		cobra.CheckErr(err)

// 		// Temp Hack, return the Cue instead
// 		cueString := ""

// 		for fields.Next() {
// 			if fields.IsDefinition() {
// 				fmt.Println("Found a #Definition called", fields.Label())
// 				c := fields.Value()
// 				fmt.Println(c)

// 				// cueString += c
// 			}
// 		}

// 		fmt.Println("New MOdel with Cue: ", cueString)
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(schemaCmd)
// }
