package proxy

import (
	"fmt"

	"github.com/ybbus/jsonrpc"
)

// ErrProxy is for general errors that originate inside the proxy module
const ErrProxy int = -32080

// ErrInternal is a general server error code
const ErrInternal int = -32603

// ErrAuthFailed is when supplied auth_token / account_id is not present in the database.
const ErrAuthFailed int = -32085

// ErrJSONParse means invalid JSON was received by the server.
const ErrJSONParse int = -32700

// ErrInvalidParams signifies a client-supplied params error
const ErrInvalidParams int = -32602

// ErrInvalidRequest signifies a general client error
const ErrInvalidRequest int = -32600

// ErrMethodUnavailable means the client-requested method cannot be found
const ErrMethodUnavailable int = -32601

// CallError is for whatever errors might occur when processing or forwarding client JSON-RPC request
type CallError interface {
	AsRPCResponse() *jsonrpc.RPCResponse
	Code() int
	Error() string
}

type GenericError struct {
	err  error
	code int
}

// InputError is a client JSON parsing error
type InputError struct {
	GenericError
}

// AuthFailed is for authentication failures when jsonrpc client has provided a token
type AuthFailed struct {
	err error
}

// AsRPCResponse returns error as jsonrpc.RPCResponse
func (e GenericError) AsRPCResponse() *jsonrpc.RPCResponse {
	return &jsonrpc.RPCResponse{
		Error: &jsonrpc.RPCError{
			Code:    e.Code(),
			Message: e.Error(),
		},
		JSONRPC: "2.0",
	}
}

// NewError is for general internal errors
func NewError(e error) GenericError {
	return GenericError{e, ErrInternal}
}

// NewInputError is for client JSON parsing errors
func NewInputError(e error) InputError {
	return InputError{GenericError{e, ErrJSONParse}}
}

// NewMethodError creates a call method error
func NewMethodError(e error) GenericError {
	return GenericError{e, ErrMethodUnavailable}
}

// NewParamsError signifies an error in method parameters
func NewParamsError(e error) GenericError {
	return GenericError{e, ErrInvalidParams}
}

// NewInternalError is for SDK-related errors (connection problems etc)
func NewInternalError(e error) GenericError {
	return GenericError{e, ErrInternal}
}

func (e GenericError) Error() string {
	return e.err.Error()
}

// Code returns JSRON-RPC error code
func (e GenericError) Code() int {
	return e.code
}

func (e GenericError) Unwrap() error {
	return e.err
}

func (e AuthFailed) Error() string {
	return fmt.Sprintf("couldn't find account for in lbrynet")
}

// Code returns JSRON-RPC error code
func (e AuthFailed) Code() int {
	return ErrAuthFailed
}

// Code returns JSRON-RPC error code
func (e InputError) Code() int {
	return ErrJSONParse
}
