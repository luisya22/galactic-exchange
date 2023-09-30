package corporation

type CrewMember struct {
	ID         uint64
	Name       string
	Species    string
	Skills     map[string]int
	AssignedTo uint64
}
