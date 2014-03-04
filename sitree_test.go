package sitree

import (
    "fmt"
    "testing"
    "time"
)

//------------------------------------------------------------
// Tests
//------------------------------------------------------------

func TestGenerate(t *testing.T) {
    fmt.Println("\n\n\n\n")

    st := New("domain.com")

    w1 := &WalkerImpl{
        name: "Branch A",
        pages: []*Leaf{
            &Leaf{Loc: "/a", Mobile: true},
            &Leaf{Loc: "/a/a", Priority: 0.8, Mobile: true},
            &Leaf{Loc: "/a/b", Lastmod: time.Now(), Changefreq: "monthly", Priority: 0.8},
            },
        }
    w2 := &WalkerImpl{
        name: "Branch B",
        pages: []*Leaf{
            &Leaf{Loc: "/b"},
            &Leaf{Loc: "/b/a"},
            &Leaf{Loc: "/b/b"},
            },
        }

    st.AddBranch(w1)
    st.AddBranch(w2)
    fmt.Println(st)
    fmt.Println("Size =", st.Size())

    // Output some XML
    st.GenerateSitemap("test_sitemap.xml")

    // Set to periodic execution
    sched := Scheduler{}
    err := sched.Set(0, 1, generator)
    if err != nil {
        t.Errorf("Error while scheduling periodic sitemap generation: %s", err)
    }

     // Forever loop to test reloading
     fmt.Println("Running continuously, use Ctrl-C to exit...")
     ch := make(chan bool, 1)
     <-ch

    fmt.Println("FINISHED")
}

//------------------------------------------------------------
// Example server process implementation
//------------------------------------------------------------

func generator() {
    fmt.Println("Generating domain tree...")
    fmt.Println("OK")
}

//------------------------------------------------------------
// Example Walker implementation
//------------------------------------------------------------

type WalkerImpl struct {
    name string
    pages []*Leaf
    walker *Walker
}

func (w *WalkerImpl) Walk() (string, []*Leaf) {
    return w.name, w.pages
}
