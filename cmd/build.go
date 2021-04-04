package cmd

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/devrel-blox/drb/blox"
	"github.com/devrel-blox/drb/config"
	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build content source into a JSON file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		cobra.CheckErr(err)

		// convert markdown to yaml
		cobra.CheckErr(convertModels(cfg))

		// validate yaml
		cobra.CheckErr(validateModels(cfg))

		data, err := aggregateModels(cfg)
		cobra.CheckErr(err)

		f, err := os.Create(path.Join(cfg.Base, cfg.Destination, "data.json"))
		cobra.CheckErr(err)

		defer f.Close()

		bb, err := json.Marshal(data)
		cobra.CheckErr(err)

		_, err = f.Write(bb)

		cobra.CheckErr(err)
		pterm.Info.Printf("wrote: %s\n", path.Join(cfg.Base, cfg.Destination, "data.json"))

	},
}

func aggregateModels(cfg *config.BloxConfig) (map[string][]interface{}, error) {
	var errors error

	data := make(map[string][]interface{})
	pterm.Info.Println("Aggregating data files...")

	for _, model := range blox.Models {
		// Attempt to decode all the YAML files with this directory as model

		filepath.Walk(path.Join(cfg.Base, cfg.Destination, model.Folder),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					// Squash, we've not even validated that it's a supported ext
					return nil
				}

				if info.IsDir() {
					return nil
				}

				ext := filepath.Ext(path)

				// if ext != cfg.DefaultExtension {
				// Should be SupportedExtensions?
				if ext != ".yaml" && ext != ".yml" {
					return nil
				}

				cueSchema := model.Cue
				if replace, ok := cfg.SchemaOverrides.Replace[model.ID]; ok {
					cueSchema = replace
				}

				entity, err := blox.FromYAML(path, model.ID, cueSchema)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, path))
					return err
				}
				data[model.Folder] = append(data[model.ID], entity)

				return nil

			})
	}
	if errors != nil {
		pterm.Error.Println("Aggregation failed")
	} else {
		pterm.Success.Println("Aggregation complete")
	}
	return data, nil
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
