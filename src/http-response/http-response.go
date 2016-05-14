package httpResponse

import "github.com/gin-gonic/gin"

// Resource represents a unique data
type Resource struct {
	Result interface{} `json:"result"`
}

// Collection represents a list of data
type Collection struct {
	Results []interface{} `json:"results"`
}

type errorReason struct {
	Reason string `json:"reason"`
}

// CodeOk contains the HTTP code for when a GET/PUT request succeeded
const CodeOk = 200

// Ok returns a CodeOk json response
func Ok(gin *gin.Context, data interface{}) {
	gin.JSON(CodeOk, data)
}

// CodeCreated contains the HTTP code for when the POST (create) request succeeded
const CodeCreated = 201

// Created returns a CodeCreated json response
func Created(gin *gin.Context, data interface{}) {
	gin.JSON(CodeCreated, data)
}

// CodeNoContent contains the HTTP code for when the request succeeded and does not returns any data
const CodeNoContent = 204

// NoContent returns a CodeNoContent json response
func NoContent(gin *gin.Context) {
	gin.Writer.WriteHeader(CodeNoContent)
}

// CodeBadRequest contains the HTTP code for when the request fails because of a user provided data
const CodeBadRequest = 400

// BadRequest returns a CodeBadRequest json response
func BadRequest(gin *gin.Context, reason string) {
	gin.JSON(CodeBadRequest, errorReason{reason})
}

// CodeUnauthorized contains the HTTP code for when the request fails because the user has not logged in
const CodeUnauthorized = 401

// Unauthorized returns a CodeUnauthorized json response
func Unauthorized(gin *gin.Context) {
	gin.Writer.WriteHeader(CodeUnauthorized)
}

// CodeForbiden contains the HTTP code for when the request fails because the user is logged in and has not the permission
const CodeForbiden = 403

// Forbiden returns a CodeForbiden json response
func Forbiden(gin *gin.Context) {
	gin.Writer.WriteHeader(CodeForbiden)
}

// CodeNotFound contains the HTTP code for when the request fails because the requested content does not exists
const CodeNotFound = 404

// NotFound returns a CodeNotFound json response
func NotFound(gin *gin.Context) {
	gin.Writer.WriteHeader(CodeNotFound)
}

// CodeConflict contains the HTTP code for when the request fails because there's a conflict with the request (duplicate data, etc)
const CodeConflict = 409

// Conflict returns a CodeConflict json response
func Conflict(gin *gin.Context, reason string) {
	gin.JSON(CodeConflict, errorReason{reason})
}

// CodeServerError contains the HTTP code for when the request fails because of a server error
const CodeServerError = 500

// ServerError returns a CodeServerError json response
func ServerError(gin *gin.Context) {
	gin.Writer.WriteHeader(CodeServerError)
}
