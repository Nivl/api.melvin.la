package education

import (
	"errors"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/datetime"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

var updateEndpoint = &router.Endpoint{
	Verb:    "PATCH",
	Path:    "/about/education/{id}",
	Handler: Update,
	Guard: &guard.Guard{
		ParamStruct: &UpdateParams{},
		Auth:        guard.AdminAccess,
	},
}

// UpdateParams represents the request params accepted by the Update handler
type UpdateParams struct {
	ID             string  `from:"url" json:"id" params:"required,uuid"`
	OrganizationID *string `from:"form" json:"organization_id" params:"noempty,uuid" maxlen:"255"`
	Degree         *string `from:"form" json:"degree" params:"noempty,trim" maxlen:"255"`
	GPA            *string `from:"form" json:"gpa" params:"trim" maxlen:"5"`
	Location       *string `from:"form" json:"location" params:"noempty,trim" maxlen:"255"`
	Description    *string `from:"form" json:"description" params:"noempty,trim" maxlen:"10000"`
	StartYear      *int    `from:"form" json:"start_year"`
	EndYear        *int    `from:"form" json:"end_year"`
	InTrash        *bool   `from:"form" json:"in_trash"`
	UnsetEndYear   bool    `from:"form" json:"unset_end_year" default:"false"`
}

// IsValid implements the params.CustomValidation interface
func (p *UpdateParams) IsValid() (isValid bool, fieldFailing string, err error) {
	if p.StartYear != nil {
		if *p.StartYear < 1900 || *p.StartYear > 2100 {
			return false, "start_year", errors.New(ErrMsgInvalidStartYear)
		}
	}

	if p.EndYear != nil {
		if *p.EndYear < 1900 || *p.EndYear > 2100 {
			return false, "end_year", errors.New(ErrMsgInvalidEndYear)
		}
	}

	if p.StartYear != nil && p.EndYear != nil {
		if *p.EndYear < *p.StartYear {
			return false, "end_year", errors.New(ErrMsgEndYearBeforeStart)
		}
	}
	return true, "", nil
}

// Update represent an API handler to update an organization
func Update(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*UpdateParams)

	edu, err := GetAnyByID(deps.DB, params.ID)
	if err != nil {
		return err
	}

	if params.Degree != nil {
		edu.Degree = *params.Degree
	}
	if params.Degree != nil {
		edu.GPA = *params.GPA
	}
	if params.Location != nil {
		edu.Location = *params.Location
	}
	if params.Description != nil {
		edu.Description = *params.Description
	}
	if params.StartYear != nil {
		edu.StartYear = *params.StartYear
	}
	if params.EndYear != nil {
		edu.EndYear = params.EndYear
	}
	if params.UnsetEndYear {
		edu.EndYear = nil
	}
	if params.OrganizationID != nil {
		org, err := organizations.GetByID(deps.DB, *params.OrganizationID)
		if err != nil {
			if apierror.IsNotFound(err) {
				return apierror.NewNotFoundField("organization_id", err.Error())
			}
			return err
		}
		edu.Organization = org
		edu.OrganizationID = org.ID
	}
	if params.InTrash != nil {
		if *params.InTrash && edu.DeletedAt == nil {
			edu.DeletedAt = datetime.Now()
		} else {
			edu.DeletedAt = nil
		}
	}

	// Let's make sure the end date is not set before the start date
	if edu.EndYear != nil {
		if *edu.EndYear < edu.StartYear {
			return apierror.NewBadRequest("end_date", ErrMsgEndYearBeforeStart)
		}
	}

	if err := edu.Update(deps.DB); err != nil {
		return err
	}

	req.Response().Ok(edu.ExportPrivate())
	return nil
}
