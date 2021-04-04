package cmd

import (
	"os"
	"path"
	"path/filepath"

	"github.com/devrel-blox/drb/blox"
	"github.com/devrel-blox/drb/config"
	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
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
		cfg, err := config.Load()
		cobra.CheckErr(err)
		cobra.CheckErr(validateModels(cfg))

	},
}

func validateModels(cfg *config.BloxConfig) error {
	var errors error

	// Create and start a fork of the default spinner.
	pterm.Info.Println("Validating YAML Files...")

	// We want to validate all the YAML for the models that we're aware of.
	for _, model := range blox.Models {
		// Attempt to decode all the YAML files with this directory as model

		filepath.Walk(path.Join(cfg.Base, cfg.Destination, model.Folder),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					// Squash, we've not even validated that it's a supported ext
					return err
				}

				if info.IsDir() {
					return err
				}

				ext := filepath.Ext(path)

				// if ext != cfg.DefaultExtension {
				// Should be SupportedExtensions?
				if ext != ".yaml" && ext != ".yml" {
					return err
				}

				cueSchema := model.Cue
				if replace, ok := cfg.SchemaOverrides.Replace[model.ID]; ok {
					cueSchema = replace
				}

				_, err = blox.FromYAML(path, model.ID, cueSchema)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, path))
					return err
				}

				return err

			})
	}
	if errors != nil {
		pterm.Error.Println("Validations failed")
	} else {
		pterm.Success.Println("Validations complete")
	}

	return errors
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
