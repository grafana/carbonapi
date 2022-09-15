package errors

import (
	"fmt"
	"net/http"
)

// ErrMissingExpr is a parse error returned when an expression is missing.
type ErrMissingExpr string

func (e ErrMissingExpr) Error() string {
	return fmt.Sprintf(string(e))
}

func (e ErrMissingExpr) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrMissingComma is a parse error returned when an expression is missing a comma.
type ErrMissingComma string

func (e ErrMissingComma) Error() string {
	return fmt.Sprintf("missing comma: %s", string(e))
}

func (e ErrMissingComma) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrMissingQuote is a parse error returned when an expression is missing a quote.
type ErrMissingQuote string

func (e ErrMissingQuote) Error() string {
	return fmt.Sprintf("missing quote: %s", string(e))
}

func (e ErrMissingQuote) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrUnexpectedCharacter is a parse error returned when an expression contains an unexpected character.
type ErrUnexpectedCharacter struct {
	Expr    string
	CharNum int
	Char    string
}

func (e ErrUnexpectedCharacter) Error() string {
	return fmt.Sprintf("unexpected character. string_to_parse=%s character_number=%d character=%s", e.Expr, e.CharNum, e.Char)
}

func (e ErrUnexpectedCharacter) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrBadType is an eval error returned when an argument has the wrong type.
type ErrBadType struct {
	Exp string
	Got string
}

func (e ErrBadType) Error() string {
	return fmt.Sprintf("bad type. expected %s - got %s", e.Exp, e.Got)
}

func (e ErrBadType) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrUnknownFunction is an error that is returned when an unknown function is specified in the query
type ErrUnknownFunction string

func (e ErrUnknownFunction) Error() string {
	return fmt.Sprintf("unknown function %q", string(e))
}

func (e ErrUnknownFunction) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrMissingArgument is an eval error returned when an argument is missing.
type ErrMissingArgument struct {
	Target string
}

func (e ErrMissingArgument) Error() string {
	return fmt.Sprintf("missing argument for: %q", e.Target)
}

func (e ErrMissingArgument) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrMissingTimeseries is an eval error returned when a time series argument is missing.
type ErrMissingTimeSeries struct {
	Target string
}

func (e ErrMissingTimeSeries) Error() string {
	return fmt.Sprintf("missing time series argument for: %s", e.Target)
}

func (e ErrMissingTimeSeries) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrUnknownTimeUnits is an eval error returned when a time unit is unknown to system
type ErrUnknownTimeUnits struct {
	Units string
}

func (e ErrUnknownTimeUnits) Error() string {
	return fmt.Sprintf("unknown time units: %s", e.Units)
}

func (e ErrUnknownTimeUnits) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrUnsupportedConsolidationFunction is an eval error returned when a consolidation function is unknown to system
type ErrUnsupportedConsolidationFunction string

func (e ErrUnsupportedConsolidationFunction) Error() string {
	return fmt.Sprintf("unknown consolidation function %q", string(e))
}

func (e ErrUnsupportedConsolidationFunction) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrBadData
type ErrBadData struct {
	Target string
	Msg    string
}

func (e ErrBadData) Error() string {
	return fmt.Sprintf("%s: bad data: %s", e.Target, e.Msg)
}

func (e ErrBadData) Message() string {
	return e.Msg
}

func (e ErrBadData) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrWildcardNotAllowed is an eval error returned when a wildcard/glob argument is found where a single series is required.
type ErrWildcardNotAllowed string

func (e ErrWildcardNotAllowed) Error() string {
	return fmt.Sprintf("\"found wildcard where series expected\" %q", string(e))
}

func (e ErrWildcardNotAllowed) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrTooManyArguments is an eval error returned when too many arguments are provided.
type ErrTooManyArguments string

func (e ErrTooManyArguments) Error() string {
	return fmt.Sprintf("\"too many arguments\" %q", string(e))
}

func (e ErrTooManyArguments) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// ErrInvalidArgument
type ErrInvalidArgument string

func (e ErrInvalidArgument) Error() string {
	return fmt.Sprintf("\"too many arguments\" %q", string(e))
}

func (e ErrInvalidArgument) HTTPStatusCode() int {
	return http.StatusBadRequest
}
