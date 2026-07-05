package variables

import "strings"

// Values holds the variables that will be replaced in templates.
type Values struct {
	ProjectName string
	ModuleName  string
	Year        string
	Author      string
}

// Replace applies variable substitutions to the input string.
func (v *Values) Replace(input string) string {
	r := strings.NewReplacer(
		"{{PROJECT_NAME}}", v.ProjectName,
		"{{MODULE_NAME}}", v.ModuleName,
		"{{YEAR}}", v.Year,
		"{{AUTHOR}}", v.Author,
	)
	return r.Replace(input)
}
