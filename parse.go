package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/erizocosmico/lua/ast"
)

const (
	// Driver name.
	Driver = "lua:1.0.0"
	// Language name.
	Language = "lua"
	// Version of the driver.
	Version = "1.0.0"

	// ParseASTAction is the default action of the driver.
	ParseASTAction = "ParseAST"
)

// AST is a collection of lua statements.
type AST struct {
	// Stmts contains all the statements in the lua file.
	Stmts []ast.Stmt
}

// Request input received by the driver.
type Request struct {
	// Action of the request.
	Action string `json:"action"`
	// Content of the request, which contains the source code to parse.
	Content string `json:"content"`
	// a request can have more fields, but we only care about these
}

// Response is the result output of parsing the AST.
type Response struct {
	// AST contains the AST of the received lua content.
	AST *AST `json:"ast"`
	// Status of the response.
	Status Status `json:"status"`
	// Errors occurred during the process of the request.
	Errors []ErrorMessage `json:"errors"`
	// Language of the driver.
	Language string `json:"language"`
	// Version of the current driver.
	Version string `json:"language_version"`
	// Driver identifier.
	Driver string `json:"driver"`
}

// Status of a response.
type Status string

const (
	// Ok means the result was correct and the AST could be parsed.
	Ok Status = "ok"
	// Error means a non fatal error happened, either because the request had
	// incorrect data or the parse was not correct.
	Error Status = "error"
	// Fatal means a fatal error happened and the request could not be handled.
	Fatal Status = "fatal"
)

// DefaultResponse creates a new empty response with the driver settings set.
func DefaultResponse() *Response {
	return &Response{
		Language: Language,
		Driver:   Driver,
		Version:  Version,
		Errors:   []ErrorMessage{},
	}
}

// NewASTResponse creates a new response with an AST.
func NewASTResponse(ast *AST) *Response {
	resp := DefaultResponse()
	resp.Status = Ok
	resp.AST = ast
	return resp
}

// NewErrorResponse creates a new response with an error message.
func NewErrorResponse(status Status, err error) *Response {
	resp := DefaultResponse()
	resp.Status = status
	resp.Errors = append(resp.Errors, ErrorMessage{status, fmt.Sprint(err)})
	return resp
}

// ErrorMessage is a single error that occurred during the process of the
// request.
type ErrorMessage struct {
	// Level of severity of the error.
	Level Status `json:"level"`
	// Message of the error.
	Message string `json:"message"`
}

func processRequest(line []byte) *Response {
	var req Request
	if err := json.Unmarshal(line, &req); err != nil {
		return NewErrorResponse(
			Fatal,
			fmt.Errorf("unable to decode request from json: %s", err),
		)
	}

	if req.Action != ParseASTAction {
		return NewErrorResponse(
			Error,
			fmt.Errorf("unknown action: %s", req.Action),
		)
	}

	a := new(AST)
	var err error
	a.Stmts, err = ast.Parse(req.Content, 1)
	if err != nil {
		return NewErrorResponse(
			Error,
			fmt.Errorf("error parsing content: %s", err),
		)
	}

	return NewASTResponse(a)
}

var pretty = flag.Bool("pretty", false, "pretty print output json")

func main() {
	flag.Parse()
	var in = bufio.NewReader(os.Stdin)

	for {
		line, _, err := in.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			// irrecoverable error, there's something wrong with the reader
			os.Exit(-1)
		}

		out, err := marshal(processRequest(line))
		if err != nil {
			out, err = marshal(NewErrorResponse(
				Fatal,
				fmt.Errorf("unable to encode to json: %s", err),
			))
			if err != nil {
				os.Exit(-1)
			}
		}

		fmt.Println(string(out))
	}
}

func marshal(resp *Response) ([]byte, error) {
	if *pretty {
		return json.MarshalIndent(resp, "", "  ")
	}
	return json.Marshal(resp)
}
