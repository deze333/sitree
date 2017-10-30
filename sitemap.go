package sitree

import (
    "bytes"
    "io"
    "os"
    "text/template"
)

//------------------------------------------------------------
// Template
//------------------------------------------------------------

var _tplTxt = `
<?xml version="1.0" encoding="UTF-8" ?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
  xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0">
{{range .Branches}}
    {{range .Leafs}}
    <url>
        <loc>{{.Loc}}</loc>
        {{if .Mobile}}
        <mobile:mobile/>
        {{end}}
        {{if .LastmodTxt}}
        <lastmod>{{.LastmodTxt}}</lastmod>
        {{end}}
        {{if .Changefreq}}
        <changefreq>{{.Changefreq}}</changefreq>
        {{end}}
        {{if .Priority}}
        <priority>{{.Priority}}</priority>
        {{end}}
    </url>
    {{end}}
{{end}}
</urlset>
`

var _tpl = template.Must(template.New("sitemap").Parse(_tplTxt))

// Generates sitemap file. First creates a temporary file,
// then replaces requested file with temporary.
func (t *Tree) GenerateSitemap(fpath string) (err error) {
	/* Allow empty sitemaps
    if t.Size() == 0 {
        return errors.New("Sitemap is empty, generation skipped")
    }
	*/
    
    wr := &SitemapWriter{}
    if err = wr.CreateFile(fpath); err != nil {
        return
    }

    // Execute template
    err = _tpl.ExecuteTemplate(wr, "sitemap", t)

    // Close file
    if err = wr.CloseFile(); err != nil {
        return
    }

    // Delete possibly existing file
    if _, err = os.Stat(fpath); err == nil {
        if err = os.Remove(fpath); err != nil {
            return
        }
    }
    // Rename temporary file
    err = os.Rename(fpath + ".tmp", fpath)
    return
}

//------------------------------------------------------------
// Writer
//------------------------------------------------------------

// Writer that attempts to remove empty lines left by template execution.
type SitemapWriter struct {
    file   *os.File
    writer *io.Writer
}

//------------------------------------------------------------
// Writer methods
//------------------------------------------------------------

func (sw *SitemapWriter) CreateFile(fpath string) (err error) {
    sw.file, err = os.Create(fpath + ".tmp")
    return
}

func (sw *SitemapWriter) CloseFile() (err error) {
    err = sw.file.Close()
    return
}

func (sw *SitemapWriter) Write(p []byte) (n int, err error) {
    plen := len(p)

    // Contains any chars or empty?
    var notEmpty bool
    var newlines = 0
    buf := bytes.NewBuffer(p)
    for {
        r, _, err := buf.ReadRune()
        if err == io.EOF {
            break
        }
        if r == '\n' {
            newlines++
        }
        if r != ' ' && r != '\n' {
            notEmpty = true
        }
    }

    if notEmpty {
        // Not empty string, allow max 1 end newline
        p = bytes.TrimLeft(p, "\n")
        p = bytes.TrimRight(p, " ")
    } else {
        // Empty string
        if newlines > 0 {
            p = []byte{}
        } 
    }

    // Write to file
    _, err = sw.file.Write(p)
    return plen, err
}
