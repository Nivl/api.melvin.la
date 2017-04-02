package generate

var modelTpl = `package {{.PackageName}}

// Code auto-generated; DO NOT EDIT

import (
	"errors"
	{{ if .Generate "JoinSQL" -}}
	"fmt"
	{{- end }}

	"github.com/melvin-laplanche/ml-api/src/core/network/http/httperr"
	"github.com/melvin-laplanche/ml-api/src/core/storage/db"
	uuid "github.com/satori/go.uuid"
)

{{ if .Generate "JoinSQL" -}}
// Join{{.ModelName}}SQL returns a string ready to be embed in a JOIN query
func Join{{.ModelName}}SQL(prefix string) string {
	fields := []string{ {{.FieldsAsArray}} }
	output := ""

	for i, field := range fields {
		if i != 0 {
			output += ", "
		}

		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\"", fullName, fullName)
	}

	return output
}
{{- end }}

{{ if .Generate "Get" -}}
// Get{{.ModelName}} finds and returns an active {{.ModelNameLC}} by ID
func Get{{.ModelName}}(id string) (*{{.ModelName}}, error) {
	{{.ModelVar}} := &{{.ModelName}}{}
	stmt := "SELECT * from {{.TableName}} WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := db.Get({{.ModelVar}}, stmt, id)
	// We want to return nil if a {{.ModelNameLC}} is not found
	if {{.ModelVar}}.ID == "" {
		return nil, err
	}
	return {{.ModelVar}}, err
}
{{- end }}

{{ if .Generate "Exists" -}}
// {{.ModelName}}Exists checks if a {{.ModelNameLC}} exists for a specific ID
func {{.ModelName}}Exists(id string) (bool, error) {
	exists := false
	stmt := "SELECT exists(SELECT 1 FROM {{.TableName}} WHERE id=$1 and deleted_at IS NULL)"
	err := db.Writer.Get(&exists, stmt, id)
	return exists, err
}
{{- end }}

{{ if .Generate "Save" -}}
// Save creates or updates the {{.ModelNameLC}} depending on the value of the id
func ({{.ModelVar}} *{{.ModelName}}) Save() error {
	return {{.ModelVar}}.SaveQ(db.Writer)
}
{{- end }}

{{ if .Generate "SaveQ" -}}
// SaveQ creates or updates the article depending on the value of the id using
// a transaction
func ({{.ModelVar}} *{{.ModelName}}) SaveQ(q db.Queryable) error {
	if {{.ModelVar}} == nil {
		return httperr.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return {{.ModelVar}}.CreateQ(q)
	}

	return {{.ModelVar}}.UpdateQ(q)
}
{{- end }}

{{ if .Generate "Create" -}}
// Create persists a {{.ModelNameLC}} in the database
func ({{.ModelVar}} *{{.ModelName}}) Create() error {
	return {{.ModelVar}}.CreateQ(db.Writer)
}
{{- end }}

{{ if .Generate "CreateQ" -}}
// Create persists a {{.ModelNameLC}} in the database
func ({{.ModelVar}} *{{.ModelName}}) CreateQ(q db.Queryable) error {
	if {{.ModelVar}} == nil {
		return httperr.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID != "" {
		return httperr.NewServerError("cannot persist a {{.ModelNameLC}} that already has an ID")
	}

	return {{.ModelVar}}.doCreate(q)
}
{{- end }}

{{ if .Generate "doCreate" -}}
// doCreate persists a {{.ModelNameLC}} in the database using a Node
func ({{.ModelVar}} *{{.ModelName}}) doCreate(q db.Queryable) error {
	if {{.ModelVar}} == nil {
		return errors.New("{{.ModelNameLC}} not instanced")
	}

	{{.ModelVar}}.ID = uuid.NewV4().String()
	{{.ModelVar}}.UpdatedAt = db.Now()
	if {{.ModelVar}}.CreatedAt == nil {
		{{.ModelVar}}.CreatedAt = db.Now()
	}

	stmt := "{{.CreateStmt}}"
	_, err := q.NamedExec(stmt, {{.ModelVar}})

  return err
}
{{- end }}

{{ if .Generate "Update" -}}
// Update updates most of the fields of a persisted {{.ModelNameLC}}.
// Excluded fields are id, created_at, deleted_at, etc.
func ({{.ModelVar}} *{{.ModelName}}) Update() error {
	return {{.ModelVar}}.UpdateQ(db.Writer)
}
{{- end }}

{{ if .Generate "UpdateQ" -}}
// Update updates most of the fields of a persisted {{.ModelNameLC}} using a transaction
// Excluded fields are id, created_at, deleted_at, etc.
func ({{.ModelVar}} *{{.ModelName}}) UpdateQ(q db.Queryable) error {
	if {{.ModelVar}} == nil {
		return httperr.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted {{.ModelNameLC}}")
	}

	return {{.ModelVar}}.doUpdate(q)
}
{{- end }}

{{ if .Generate "doUpdate" -}}
// doUpdate updates a {{.ModelNameLC}} in the database using an optional transaction
func ({{.ModelVar}} *{{.ModelName}}) doUpdate(q db.Queryable) error {
	if {{.ModelVar}} == nil {
		return httperr.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted {{.ModelNameLC}}")
	}

	{{.ModelVar}}.UpdatedAt = db.Now()

	stmt := "{{.UpdateStmt}}"
	_, err := q.NamedExec(stmt, {{.ModelVar}})

	return err
}
{{- end }}

{{ if .Generate "FullyDelete" -}}
// FullyDelete removes a {{.ModelNameLC}} from the database
func ({{.ModelVar}} *{{.ModelName}}) FullyDelete() error {
	return {{.ModelVar}}.FullyDeleteQ(db.Writer)
}
{{- end }}

{{ if .Generate "FullyDeleteQ" -}}
// FullyDeleteQ removes a {{.ModelNameLC}} from the database using a transaction
func ({{.ModelVar}} *{{.ModelName}}) FullyDeleteQ(q db.Queryable) error {
	if {{.ModelVar}} == nil {
		return errors.New("{{.ModelNameLC}} not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return errors.New("{{.ModelNameLC}} has not been saved")
	}

	stmt := "DELETE FROM {{.TableName}} WHERE id=$1"
	_, err := q.Exec(stmt, {{.ModelVar}}.ID)

	return err
}
{{- end }}

{{ if .Generate "Delete" -}}
// Delete soft delete a {{.ModelNameLC}}.
func ({{.ModelVar}} *{{.ModelName}}) Delete() error {
	return {{.ModelVar}}.DeleteQ(db.Writer)
}
{{- end }}

{{ if .Generate "DeleteQ" -}}
// DeleteQ soft delete a {{.ModelNameLC}} using a transaction
func ({{.ModelVar}} *{{.ModelName}}) DeleteQ(q db.Queryable) error {
	return {{.ModelVar}}.doDelete(q)
}
{{- end }}

{{ if .Generate "doDelete" -}}
// doDelete performs a soft delete operation on a {{.ModelNameLC}} using an optional transaction
func ({{.ModelVar}} *{{.ModelName}}) doDelete(q db.Queryable) error {
	if {{.ModelVar}} == nil {
		return httperr.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return httperr.NewServerError("cannot delete a non-persisted {{.ModelNameLC}}")
	}

	{{.ModelVar}}.DeletedAt = db.Now()

	stmt := "UPDATE {{.TableName}} SET deleted_at = $2 WHERE id=$1"
	_, err := q.Exec(stmt, {{.ModelVar}}.ID, {{.ModelVar}}.DeletedAt)
	return err
}
{{- end }}

{{ if .Generate "IsZero" -}}
// IsZero checks if the object is either nil or don't have an ID
func ({{.ModelVar}} *{{.ModelName}}) IsZero() bool {
	return {{.ModelVar}} == nil || {{.ModelVar}}.ID == ""
}
{{- end }}`
