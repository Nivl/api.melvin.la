package experience

import (
	"errors"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-types/date"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
)

var addEndpoint = &router.Endpoint{
	Verb:    "POST",
	Path:    "/about/experience",
	Handler: Add,
	Guard: &guard.Guard{
		Auth:        guard.AdminAccess,
		ParamStruct: &AddParams{},
	},
}

// ErrMsgInvalidEndDate represents the error messages returned when end_date
// contains a date that is before start_date
const ErrMsgInvalidEndDate = "cannot be before start_date"

// AddParams represents the params accepted by the Add endpoint
type AddParams struct {
	OrganizationID string     `from:"form" json:"organization_id" params:"required,uuid" maxlen:"255"`
	JobTitle       string     `from:"form" json:"job_title" params:"required,trim" maxlen:"255"`
	Location       string     `from:"form" json:"location" params:"required,trim" maxlen:"255"`
	Description    string     `from:"form" json:"description" params:"required,trim" maxlen:"10000"`
	StartDate      *date.Date `from:"form" json:"start_date" params:"required"`
	EndDate        *date.Date `from:"form" json:"end_date"`
}

// IsValid implements the params.CustomValidation interface
func (p *AddParams) IsValid() (isValid bool, fieldFailing string, err error) {
	if p.EndDate != nil {
		if p.EndDate.IsBefore(p.StartDate) {
			return false, "end_date", errors.New(ErrMsgInvalidEndDate)
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

	exp := &Experience{
		OrganizationID: org.ID,
		JobTitle:       params.JobTitle,
		Location:       params.Location,
		Description:    params.Description,
		StartDate:      params.StartDate,
		EndDate:        params.EndDate,

		Organization: org,
	}

	if err := exp.Create(deps.DB); err != nil {
		return err
	}
	return req.Response().Created(exp.ExportPrivate())
}
