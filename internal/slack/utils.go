package slack

import (
	"bytes"
	"encoding/json"
	"io"
)

// Read implements PostMessageParameters as a reader from which read the json
// encoding of the structure.
//
// When the reader reaches EOF, it is reset so it can be encoded again.
func (parameters *PostMessageParameters) Read(p []byte) (int, error) {
	if parameters.buffer != nil {
		i, err := parameters.buffer.Read(p)
		if err == io.EOF {
			parameters.buffer = nil
		}
		return i, err
	}

	parameters.buffer = &bytes.Buffer{}
	encoder := json.NewEncoder(parameters.buffer)
	if err := encoder.Encode(parameters); err != nil {
		return 0, err
	}
	return parameters.Read(p)
}
