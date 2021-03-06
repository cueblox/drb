package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/devrel-blox/drb/blox"
	"github.com/devrel-blox/drb/config"
	"github.com/devrel-blox/drb/cuedb"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	referentialIntegrity bool
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		cobra.CheckErr(err)

		database, err := cuedb.NewDatabase()
		cobra.CheckErr(err)

		// Load Schemas!
		cobra.CheckErr(database.RegisterTables(blox.ProfileCue))

		cobra.CheckErr(buildModels(cfg, &database))

		if referentialIntegrity {
			pterm.Info.Println("Checking Referential Integrity")
			err = database.ReferentialIntegrity()
			if err != nil {
				pterm.Error.Println(err)
			} else {
				pterm.Success.Println("Foreign Keys Validated")
			}
		}

		jso, err := database.MarshalJSON()
		cobra.CheckErr(err)

		fmt.Println("I should write this to a file")
		fmt.Println(string(jso))

	},
}

func buildModels(cfg *config.BloxConfig, db *cuedb.Database) error {
	var errors error

	pterm.Info.Println("Validating ...")

	for _, table := range db.GetTables() {
		err := filepath.Walk(path.Join(cfg.Base, cfg.Destination, table.Directory()),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
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

				slug := strings.Replace(filepath.Base(path), ext, "", -1)

				bytes, err := ioutil.ReadFile(path)
				if err != nil {
					return multierror.Append(err)
				}

				var istruct = make(map[string]interface{})

				err = yaml.Unmarshal(bytes, &istruct)

				if err != nil {
					return multierror.Append(err)
				}

				record := make(map[string]interface{})
				record[slug] = istruct

				err = db.Insert(table, record)
				if err != nil {
					return multierror.Append(err)
				}

				return err

			})

		if err != nil {
			errors = multierror.Append(err)
		}
	}

	if errors != nil {
		pterm.Error.Println("Validations failed")
	} else {
		pterm.Success.Println("Validations complete")
	}

	return errors
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVarP(&referentialIntegrity, "referential-integrity", "i", false, "Enforce referential integrity")
}
