package sitree

import (
	"fmt"
	"strings"
	"time"
)

//------------------------------------------------------------
// Model
//------------------------------------------------------------

type Tree struct {
	Name             string
	Created          time.Time
	CreationDuration time.Duration
	Branches         []Branch
}

type Branch struct {
	Name  string // for information only, not used in sitemap
	Leafs []*Leaf
}

type Leaf struct {
	Loc        string
	Lastmod    time.Time
	LastmodTxt string
	Changefreq string
	Priority   float64
	Mobile     bool
}

type Walker interface {
	Walk() (string, []*Leaf)
}

//------------------------------------------------------------
// Constructor
//------------------------------------------------------------

func New(name string) (t *Tree) {
	t = &Tree{Name: name}
	return
}

//------------------------------------------------------------
// Methods
//------------------------------------------------------------

func (t *Tree) String() string {
	tree := []string{}
	for i, branch := range t.Branches {
		b := fmt.Sprint("\n###### ", i, " : ", branch.Name, "\n")
		ls := []string{b}
		for j, leaf := range branch.Leafs {
			lss := []string{}
			lss = append(lss, fmt.Sprint(i, ".", j, " : Loc = ", leaf.Loc))
			if leaf.LastmodTxt != "" {
				lss = append(lss, fmt.Sprint("Lastmod = ", leaf.LastmodTxt))
			}

			if leaf.Changefreq != "" {
				lss = append(lss, fmt.Sprint("Changefreq = ", leaf.Changefreq))
			}
			if leaf.Priority != 0 {
				lss = append(lss, fmt.Sprint("Priority = ", leaf.Priority))
			}
			ls = append(ls, strings.Join(lss, "\n"))
		}
		tree = append(tree, strings.Join(ls, "\n------\n"))
	}
	return strings.Join(tree, "\n")
}

func (t *Tree) Size() int {
	size := 0
	for _, branch := range t.Branches {
		size += len(branch.Leafs)
	}
	return size
}

// Flushes (empties) tree so it's ready for next time.
func (t *Tree) Flush() {
    t.Branches = []Branch{}
}

func (t *Tree) AddBranches(ws []Walker) {
	for _, w := range ws {
		t.AddBranch(w)
	}
}

func (t *Tree) AddBranch(w Walker) {
	name, leafs := w.Walk()
	// Format leafs
	for _, l := range leafs {
		if !l.Lastmod.IsZero() {
			l.LastmodTxt = l.Lastmod.Format(time.RFC3339)
		}
	}
	b := Branch{
		Name:  name,
		Leafs: leafs,
	}
	t.Branches = append(t.Branches, b)
}
