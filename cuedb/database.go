package cuedb

import (
	"errors"
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/hashicorp/go-multierror"
)

// Database is the "world" struct. We can "insert" records
// into it and know immediately if they're valid or not.
type Database struct {
	runtime *cue.Runtime
	db      cue.Value
	tables  map[string]Table
}

func NewDatabase() (Database, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("", "")

	if nil != err {
		return Database{}, err
	}

	return Database{
		runtime: &cueRuntime,
		db:      cueInstance.Value(),
		tables:  make(map[string]Table),
	}, nil
}

// RegisterTables ensures that the cueString schema is a valid schema
// and parses the Cue to find Models within. Each Model is registered
// as a Table, provided the name is available.
func (d *Database) RegisterTables(cueString string) error {
	cueInstance, err := d.runtime.Compile(cueString, cueString)
	if nil != err {
		return err
	}

	cueValue := cueInstance.Value()

	// First, Unify whatever schemas the users want. We'll
	// do our best to extract whatever information from
	// it we require
	d.db = d.db.FillPath(cue.Path{}, cueValue)

	// Is the Schema valid?
	_, err = GetSchemaVersion(cueValue)
	if err != nil {
		return err
	}

	// We only have a V1 :)
	metadata, err := GetSchemaV1Metadata(cueValue)
	if err != nil {
		return err
	}

	// Find Models and register as a table
	fields, err := cueValue.Fields(cue.Definitions(true))
	if err != nil {
		return err
	}

	for fields.Next() {
		if !fields.IsDefinition() {
			// Only Definitions can be registered as tables
			continue
		}

		// We have a Definition, does it define a model?
		plural, err := GetV1Model(fields.Value())
		if nil != err {
			continue
		}

		table := Table{
			schemaNamespace: metadata.Namespace,
			schemaName:      metadata.Name,
			name:            fields.Label(),
			plural:          plural,
		}

		if _, ok := d.tables[table.ID()]; ok {
			return errors.New(fmt.Sprintf("Table with name '%s' already registered", fields.Label()))
		}

		err = d.db.Validate()
		if nil != err {
			return err
		}

		inst, err := d.runtime.Compile("", fmt.Sprintf("{%s: _\n%s: [ID=string]: %s}", fields.Label(), plural, fields.Label()))
		if err != nil {
			return err
		}

		d.db = d.db.FillPath(cue.Path{}, inst.Value())

		if err := d.db.Validate(); nil != err {
			return err
		}

		d.tables[table.ID()] = table
	}

	return nil
}

func (d *Database) GetTables() map[string]Table {
	return d.tables
}

func (d *Database) GetTable(name string) (Table, error) {
	if table, ok := d.tables[name]; ok {
		return table, nil
	}

	return Table{}, errors.New(fmt.Sprintf("Table '%s' doesn't exist in database", name))
}

func (d *Database) MarshalJSON() ([]byte, error) {
	return d.db.MarshalJSON()
}

type Table struct {
	name   string
	plural string // The directory where we find the records

	// Which schema registered this table?
	schemaNamespace string
	schemaName      string
}

func (t *Table) ID() string {
	return strings.ToLower(t.name)
}

func (t *Table) Directory() string {
	return t.plural
}

func (t *Table) CuePath() cue.Path {
	return cue.ParsePath(t.plural)
}

func (d *Database) Insert(table Table, record map[string]interface{}) error {
	filledValued := d.db.FillPath(table.CuePath(), record)

	err := filledValued.Validate()
	if nil != err {
		return err
	}

	d.db = d.db.Unify(filledValued)

	return nil
}

func (d *Database) ReferentialIntegrity() error {
	for _, table := range d.GetTables() {
		// Walk each field and look for _id labels
		// fmt.Println("Finding Def: ", table.name)
		val := d.db.LookupDef(table.name)

		fields, err := val.Fields(cue.Optional(true))
		if err != nil {
			return err
		}

		for fields.Next() {
			if strings.HasSuffix(fields.Label(), "_id") {
				foreignTable, err := d.GetTable(fmt.Sprintf("#%s", strings.TrimSuffix(fields.Label(), "_id")))
				if err != nil {
					return err
				}

				inst, err := d.runtime.Compile("", fmt.Sprintf("{%s: _\n%s: %s: or([ for k, _ in %s {k}])}", foreignTable.plural, table.name, fields.Label(), foreignTable.plural))
				if err != nil {
					return err
				}
				d.db = d.db.FillPath(cue.Path{}, inst.Value())
			}
		}
	}

	err := d.db.Validate()
	if err != nil {
		return multierror.Prefix(err, "Referential Integrity Failed")
	}

	return nil
}
