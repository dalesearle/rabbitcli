package data

type ByKey []string

func (k ByKey) Len() int           { return len(k) }
func (k ByKey) Less(i, j int) bool { return k[i] < k[j] }
func (k ByKey) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }