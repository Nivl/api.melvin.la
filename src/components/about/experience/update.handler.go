package experience

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/date"
	"github.com/Nivl/go-types/datetime"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

var updateEndpoint = &router.Endpoint{
	Verb:    http.MethodPatch,
	Path:    "/about/experience/{id}",
	Handler: Update,
	Guard: &guard.Guard{
		ParamStruct: &UpdateParams{},
		Auth:        guard.AdminAccess,
	},
}

// UpdateParams represents the request params accepted by the Update handler
type UpdateParams struct {
	ID             string     `from:"url" json:"id" params:"required,uuid"`
	OrganizationID *string    `from:"form" json:"organization_id" params:"noempty,uuid" maxlen:"255"`
	JobTitle       *string    `from:"form" json:"job_title" params:"noempty,trim" maxlen:"255"`
	Location       *string    `from:"form" json:"location" params:"noempty,trim" maxlen:"255"`
	Description    *string    `from:"form" json:"description" params:"noempty,trim" maxlen:"10000"`
	StartDate      *date.Date `from:"form" json:"start_date"`
	EndDate        *date.Date `from:"form" json:"end_date"`
	InTrash        *bool      `from:"form" json:"in_trash"`
	UnsetEndDate   bool       `from:"form" json:"unset_end_date" default:"false"`
}

// Update represent an API handler to update an organization
func Update(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*UpdateParams)

	exp, err := GetAnyByID(deps.DB, params.ID)
	if err != nil {
		return err
	}

	if params.JobTitle != nil {
		exp.JobTitle = *params.JobTitle
	}
	if params.Location != nil {
		exp.Location = params.Location
	}
	if params.Description != nil {
		exp.Description = params.Description
	}
	if params.StartDate != nil {
		exp.StartDate = params.StartDate
	}
	if params.EndDate != nil {
		exp.EndDate = params.EndDate
	}
	if params.UnsetEndDate {
		exp.EndDate = nil
	}
	if params.OrganizationID != nil {
		org, err := organizations.GetByID(deps.DB, *params.OrganizationID)
		if err != nil {
			if apierror.IsNotFound(err) {
				return apierror.NewNotFoundField("organization_id", err.Error())
			}
			return err
		}
		exp.Organization = org
		exp.OrganizationID = org.ID
	}
	if params.InTrash != nil {
		if *params.InTrash && exp.DeletedAt == nil {
			exp.DeletedAt = datetime.Now()
		} else {
			exp.DeletedAt = nil
		}
	}

	// Let's make sure the end date is not set before the start date
	if exp.EndDate != nil {
		if exp.EndDate.IsBefore(exp.StartDate) {
			return apierror.NewBadRequest("end_date", ErrMsgInvalidEndDate)
		}
	}

	if err := exp.Update(deps.DB); err != nil {
		return err
	}

	req.Response().Ok(exp.ExportPrivate())
	return nil
}
