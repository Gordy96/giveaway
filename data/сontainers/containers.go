package сontainers

import (
	"fmt"
	"giveaway/data"
	"math/rand"
	"time"
)

type Entry struct {
	Key   string
	Value interface{}
}

type EntryContainer struct {
	data  []Entry
	dupes map[string][]int
	add   func(data.HasKey)
}

func (t *EntryContainer) Get(i int) *Entry {
	if i < 0 {
		panic(fmt.Errorf("index < 0"))
	}
	return &t.data[i]
}

func (t *EntryContainer) GetRandomIndexNoDuplicates() int {
	l := len(t.data) - 1
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomizer.Seed(time.Now().UnixNano())
	var idx = -1
	for {
		idx = randomizer.Intn(l)
		entry := t.data[idx]
		if dupes := t.dupes[entry.Key]; len(dupes) == 1 && dupes[0] == idx {
			break
		}
	}
	return idx
}

func (t *EntryContainer) Add(value data.HasKey) {
	t.add(value)
}

func (t *EntryContainer) Length() int {
	return len(t.data)
}

func (t *EntryContainer) LengthNoDuplicates() int {
	return len(t.dupes)
}

func NewEntryContainer() *EntryContainer {
	c := &EntryContainer{}
	c.add = func(ins data.HasKey) {
		idx := len(c.data)

		entry := Entry{ins.GetKey().(string), ins}
		c.data = append(c.data, entry)
		if dup, in := c.dupes[entry.Key]; in {
			c.dupes[entry.Key] = append([]int{idx}, dup...)
		} else {
			c.dupes[entry.Key] = []int{idx}
		}
	}
	return c
}
