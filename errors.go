package httperror

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

/*
HttpError is error type of http.
This is maybe RFC7807(https://datatracker.ietf.org/doc/html/rfc7807).
*/
type HttpError struct {
	Type             string         `json:"type,omitempty"`
	Title            string         `json:"title,omitempty"`
	Status           int            `json:"status,omitempty"`
	Detail           string         `json:"detail,omitempty"`
	Instance         string         `json:"instance,omitempty"`
	ExtensionMembers map[string]any `json:"extension_members"`
}

// New make new HttpError instance with below inputs.
// title: A short, human-readable summary of the problem type.  It SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization (e.g., using proactive content negotiation)
// status: The HTTP status code generated by the origin server for this occurrence of the problem.
// type url: A URI reference that identifies the problem type. This specification encourages that, when dereferenced, it provide human-readable documentation for the problem type.
func New(title string, status int, typeUrl string) *HttpError {
	return &HttpError{
		Title:  title,
		Type:   typeUrl,
		Status: status,
	}
}

// WithDetail : A human-readable explanation specific to this occurrence of the problem.
func (e *HttpError) WithDetail(detail string) *HttpError {
	e.Detail = detail
	return e
}

// WithInstance : A URI reference that identifies the specific occurrence of the problem.  It may or may not yield further information if dereferenced.
func (e *HttpError) WithInstance(instance string) *HttpError {
	e.Instance = instance
	return e
}

func (e *HttpError) Error() string {
	builder := strings.Builder{}
	builder.WriteString(e.Title)
	builder.WriteString(" on ")
	builder.WriteString(e.Instance)
	builder.WriteString(" reference ")
	builder.WriteString(e.Type)

	return builder.String()
}

func (e *HttpError) ToJSON(marshal func(any) ([]byte, error)) ([]byte, error) {
	if marshal == nil {
		marshal = json.Marshal
	}

	builder := bytes.NewBuffer(nil)
	builder.WriteString("{")
	builder.WriteString("\"type\":")
	builder.WriteString("\"" + e.Type + "\"")
	builder.WriteString(",")
	builder.WriteString("\"title\":")
	builder.WriteString("\"" + e.Title + "\"")
	builder.WriteString(",")
	builder.WriteString("\"status\":")
	builder.WriteString(strconv.Itoa(e.Status))
	builder.WriteString(",")
	builder.WriteString("\"detail\":")
	builder.WriteString("\"" + e.Detail + "\"")
	builder.WriteString(",")
	builder.WriteString("\"instance\":")
	builder.WriteString("\"" + e.Instance + "\"")

	for k, v := range e.ExtensionMembers {
		builder.WriteString(k)
		builder.WriteString(":")
		d, err := marshal(v)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal extension member %s: %w", k, err)
		}

		builder.WriteString(string(d))
		builder.WriteString(",")
	}

	builder.WriteString("}")

	return builder.Bytes(), nil
}

func FromJSON(data []byte, unmarshal func([]byte, any) error) (*HttpError, error) {
	if unmarshal == nil {
		unmarshal = json.Unmarshal
	}

	v := map[string]any{}
	if err := unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("cannot unmarshal http error: %w", err)
	}

	he := &HttpError{
		Title:    v["title"].(string),
		Status:   int(v["status"].(float64)),
		Type:     v["type"].(string),
		Detail:   v["detail"].(string),
		Instance: v["instance"].(string),
	}

	delete(v, "title")
	delete(v, "status")
	delete(v, "type")
	delete(v, "detail")
	delete(v, "instance")

	he.ExtensionMembers = v

	return he, nil
}
