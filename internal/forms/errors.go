package forms

type errors map[string][]string

// Add adds an error message for a given form field
func (e errors) Add(formField, message string) {
	e[formField] = append(e[formField], message)
}

// GetFirst returns the first error message
func (e errors) GetFirst(field string) string {
	errorDesc := e[field]
	if len(errorDesc) == 0 {
		return ""
	}

	return errorDesc[0]
}
