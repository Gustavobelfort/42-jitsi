package oauth

import (
	"fmt"
	"net/url"
	"strings"
)

type Params url.Values

func (Params) toString(value interface{}) string {
	if stringer, ok := value.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%v", value)
}

// Add will append values to the parameter field, or create the field if it does not already exist.
func (p Params) Add(key string, values ...interface{}) {
	if _, ok := p[key]; !ok {
		p[key] = make([]string, 0)
	}
	for _, value := range values {
		p[key] = append(p[key], p.toString(value))
	}
}

// Set sets the value of a parameter field. If it is already filled, it overwrites its current value.
func (p Params) Set(key string, values ...interface{}) {
	p[key] = make([]string, 0)
	p.Add(key, values...)
}

// Encode encodes the parameters into a url query parameter format.
func (p Params) Encode() string {
	buffer := new(strings.Builder)
	l := len(p)
	for key, value := range p {
		buffer.WriteString(fmt.Sprintf("%s=%s", key, strings.Join(value, ",")))
		l--
		if l > 0 {
			buffer.WriteString("&")
		}
	}
	return buffer.String()
}
