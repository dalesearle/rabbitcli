package data

type ByName []*Report

func (r ByName) Len() int           { return len(r) }
func (r ByName) Less(i, j int) bool { return r[i].Name < r[j].Name }
func (r ByName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
