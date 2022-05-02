package registry

type Entry struct {
	Key      string
	SearchBy string
}

func (entry ModIndex) Len() int           { return len(entry) }
func (entry ModIndex) Less(i, j int) bool { return entry[i].SearchBy < entry[j].SearchBy }
func (entry ModIndex) Swap(i, j int)      { entry[i], entry[j] = entry[j], entry[i] }
