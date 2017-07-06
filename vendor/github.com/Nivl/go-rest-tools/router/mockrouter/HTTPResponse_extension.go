package mockrouter

import (
	"github.com/stretchr/testify/mock"
)

// ExpectCreated is a helper that expects a valid Created response
func (res *HTTPResponse) ExpectCreated(typ string, runnable func(args mock.Arguments)) *mock.Call {
	createdCall := res.On("Created", mock.AnythingOfType(typ))
	createdCall.Return(nil)
	if runnable != nil {
		createdCall.Run(runnable)
	}
	return createdCall
}

// ExpectOk is a helper that expects a valid Created response
func (res *HTTPResponse) ExpectOk(typ string, runnable func(args mock.Arguments)) *mock.Call {
	createdCall := res.On("Ok", mock.AnythingOfType(typ))
	createdCall.Return(nil)
	if runnable != nil {
		createdCall.Run(runnable)
	}
	return createdCall
}
