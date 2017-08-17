package education

import (
	"errors"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

var addEndpoint = &router.Endpoint{
	Verb:    "POST",
	Path:    "/about/education",
	Handler: Add,
	Guard: &guard.Guard{
		Auth:        guard.AdminAccess,
		ParamStruct: &AddParams{},
	},
}

// ErrMsgEndYearBeforeStart represents the error messages returned when end_year
// contains a year that is before start_year
const ErrMsgEndYearBeforeStart = "cannot be before start_year"

// ErrMsgInvalidStartYear represents the error messages returned when start_year
// is before 1900 or after 2100
const ErrMsgInvalidStartYear = "cannot be before 1900 and cannot be after 2100"

// ErrMsgInvalidEndYear represents the error messages returned when end_year
// is before 1900 or after 2100
const ErrMsgInvalidEndYear = "cannot be before 1900 and cannot be after 2100"

// AddParams represents the params accepted by the Add endpoint
type AddParams struct {
	OrganizationID string `from:"form" json:"organization_id" params:"required,uuid" maxlen:"255"`
	Degree         string `from:"form" json:"degree" params:"required,trim" maxlen:"255"`
	GPA            string `from:"form" json:"gpa" params:"trim" maxlen:"5"`
	Location       string `from:"form" json:"location" params:"required,trim" maxlen:"255"`
	Description    string `from:"form" json:"description" params:"required,trim" maxlen:"10000"`
	StartYear      int    `from:"form" json:"start_year" params:"required"`
	EndYear        *int   `from:"form" json:"end_year"`
}

// IsValid implements the params.CustomValidation interface
func (p *AddParams) IsValid() (isValid bool, fieldFailing string, err error) {
	if p.StartYear < 1900 || p.StartYear > 2100 {
		return false, "start_year", errors.New(ErrMsgInvalidStartYear)
	}

	if p.EndYear != nil {
		if *p.EndYear < 1900 || *p.EndYear > 2100 {
			return false, "end_year", errors.New(ErrMsgInvalidEndYear)
		}
		if *p.EndYear < p.StartYear {
			return false, "end_year", errors.New(ErrMsgEndYearBeforeStart)
		}
	}
	return true, "", nil
}

// Add is an endpoint used to add a new Organization
func Add(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*AddParams)

	org, err := organizations.GetByID(deps.DB, params.OrganizationID)
	if err != nil {
		return err
	}

	edu := &Education{
		OrganizationID: org.ID,
		Degree:         params.Degree,
		GPA:            params.GPA,
		Location:       params.Location,
		Description:    params.Description,
		StartYear:      params.StartYear,
		EndYear:        params.EndYear,

		Organization: org,
	}

	if err := edu.Create(deps.DB); err != nil {
		return err
	}
	return req.Response().Created(edu.ExportPrivate())
}
