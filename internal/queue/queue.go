package queue

import mod "github.com/jedwards1230/go-kerbal/internal/ckan"

type Queue struct {
	List map[string]map[string]mod.Ckan
}

func NewQueue() Queue {
	q := make(map[string]map[string]mod.Ckan, 0)

	q["remove"] = make(map[string]mod.Ckan, 0)
	q["install"] = make(map[string]mod.Ckan, 0)
	q["dependency"] = make(map[string]mod.Ckan, 0)

	return Queue{
		List: q,
	}
}

func (q *Queue) FindDependents(s string) []mod.Ckan {
	modList := make([]mod.Ckan, 0)
	for _, mod := range q.GetSelections() {
		if len(mod.ModDepends) > 0 {
			for i := range mod.ModDepends {
				if mod.ModDepends[i] == s {
					modList = append(modList, mod)
				}
			}
		}
	}
	return modList
}

func (q *Queue) CheckQueue(s string) bool {
	for _, mod := range q.GetRemovals() {
		if mod.Identifier == s {
			return true
		}
	}
	for _, mod := range q.GetSelections() {
		if mod.Identifier == s {
			return true
		}
	}
	for _, mod := range q.GetDependencies() {
		if mod.Identifier == s {
			return true
		}
	}
	return false
}

func (q Queue) CheckRemovals(s string) bool {
	for _, mod := range q.GetRemovals() {
		if mod.Identifier == s {
			return true
		}
	}
	return false
}

func (q *Queue) AddRemoval(mod mod.Ckan) {
	q.List["remove"][mod.Identifier] = mod
}

func (q *Queue) AddSelection(mod mod.Ckan) {
	q.List["install"][mod.Identifier] = mod
}

func (q *Queue) AddDependency(mod mod.Ckan) {
	q.List["dependency"][mod.Identifier] = mod
}

func (q Queue) GetRemovals() map[string]mod.Ckan {
	return q.List["remove"]
}

func (q Queue) GetSelections() map[string]mod.Ckan {
	return q.List["install"]
}

func (q Queue) GetDependencies() map[string]mod.Ckan {
	return q.List["dependency"]
}

func (q Queue) InstallLen() int {
	count := 0
	for _, mod := range q.GetSelections() {
		if !mod.Installed() {
			count += 1
		}
	}
	for _, mod := range q.GetDependencies() {
		if !mod.Installed() {
			count += 1
		}
	}
	return count
}

func (q Queue) RemoveLen() int {
	return len(q.GetRemovals())
}

func (q Queue) Len() int {
	return len(q.GetRemovals()) + len(q.GetSelections()) + len(q.GetDependencies())
}
