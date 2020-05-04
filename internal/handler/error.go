package handler

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// NoCorrectorError will be returned when the evaluation does not have any corrector.
	// This generally means that the evaluation was automatically made by the intranet.
	//
	// e.g: When a student does not make the required number of evaluations on time.
	NoCorrectorError = errors.New("the evaluation does not have any corrector")
	// NotInDBError will be returned when the evaluation is not present in the database but is expected to be.
	NotInDBError = errors.New("the evaluation was not in the database")
)

// MissingFieldsError will be returned when the evaluation's payload is missing one or multiple fields.
//
// This should result in an invalid payload and the rejection of it.
type MissingFieldsError struct {
	missing []string
}

// Error formats the MissingFieldsError dynamically depending on the missing fields list.
func (err *MissingFieldsError) Error() string {
	return fmt.Sprintf("missing required fields: %s", strings.Join(err.missing, ","))
}
