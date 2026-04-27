package role

import "strings"

const (
	ProjectAdmin   = "PROJECT_ADMIN"
	DepartmentHead = "DEPARTMENT_HEAD"
	PICProject     = "PIC_PROJECT"
	TeamMember     = "TEAM_MEMBER"
)

func FormatRoleName(raw string) string {
	acronyms := map[string]string{
		"pic": "PIC",
	}

	words := strings.Split(strings.ToLower(raw), "_")
	for i, w := range words {
		if len(w) == 0 {
			continue
		}
		if acronyms, ok := acronyms[w]; ok {
			words[i] = acronyms
		} else {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
