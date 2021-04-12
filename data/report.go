package data

type Report struct {
	Name string
	Value string
}

func NewReport(name, value string) *Report {
	return &Report{
		Name: name,
		Value: value,
	}
}

func (r *Report) TableRow () []string {
	return []string{r.Name, r.Value}
}
