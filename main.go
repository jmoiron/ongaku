package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jmoiron/ongaku/statik"
	"github.com/rakyll/statik/fs"
)

// Ongaku is a very simple http streaming media player written in Go.
// You run Ongaku in a directory with music in it, and it runs an HTTP
// server that can let you browse and play your music using the jplayer
// HTML5 player.

//go:generate -command statik statik
//go:generate statik -src=./static

const apiEndpoint = "/_api/"

var cfg struct {
	port int
}

func init() {
	flag.IntVar(&cfg.port, "port", 1339, "http port")
}

var stfs http.FileSystem
var tmpl *template.Template

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error
	stfs, err = fs.New()
	panicIf(err)

	m, err := stfs.Open("/main.tmpl")
	panicIf(err)
	src, err := ioutil.ReadAll(m)
	panicIf(err)

	tmpl = template.Must(template.New("main").Parse(string(src)))
}

func main() {
	listenOn := fmt.Sprintf(":%d", cfg.port)
	log.Printf("Listening on %s", listenOn)
	fsh := http.StripPrefix("/_file/", http.FileServer(http.Dir(".")))
	sth := http.StripPrefix("/_static/", http.FileServer(stfs))
	http.Handle("/_file/", fsh)
	http.Handle("/_static/", sth)
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc(apiEndpoint, api)
	http.HandleFunc("/", index)
	http.ListenAndServe(listenOn, nil)
}

type DirEntry struct {
	Name    string
	Path    string
	Url     string
	IsDir   bool
	IsAudio bool
}

type DirList []DirEntry

func (d DirList) Len() int      { return len(d) }
func (d DirList) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

type DirsFirst struct{ DirList }

func (d DirsFirst) Less(i, j int) bool {
	x := d.DirList
	switch {
	case x[i].IsDir && x[j].IsDir:
		return x[i].Name < x[j].Name
	case x[i].IsDir:
		return true
	case x[j].IsDir:
		return false
	default:
		return x[i].Name < x[j].Name
	}
}

// like ioutil.ReadDir, but returns a DirList
func ReadDirList(path, url string) (DirList, error) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files DirList

	for _, fi := range fis {
		path := filepath.Join(path, fi.Name())
		de := DirEntry{
			Name:    fi.Name(),
			Path:    path,
			Url:     urljoin(url, fi.Name()),
			IsDir:   fi.IsDir(),
			IsAudio: isAudio(path),
		}
		files = append(files, de)
	}

	return files, nil
}

type msi map[string]interface{}

func favicon(w http.ResponseWriter, r *http.Request) {
	fav, err := stfs.Open("favicon.ico")
	if err != nil {
		return
	}
	io.Copy(w, fav)
}

func api(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
	path := strings.TrimPrefix(r.RequestURI, apiEndpoint)
	_, abs, err := paths(path)
	if err != nil {
		logErr(w, err)
		return
	}

	files, err := ReadDirList(abs, r.RequestURI)
	if err != nil {
		logErr(w, err)
		return
	}

	var infos []AudioInfo

	for _, f := range files {
		if f.IsAudio {
			infos = append(infos, GetAudioInfo(f.Path))
		}
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(infos)
	if err != nil {
		logErr(w, err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
	rel, abs, err := paths(r.RequestURI)
	if err != nil {
		logErr(w, err)
		return
	}

	files, err := ReadDirList(abs, r.RequestURI)
	if err != nil {
		logErr(w, err)
		return
	}

	sort.Sort(DirsFirst{files})
	paths := pathSplit(rel)

	var infos []AudioInfo

	for _, f := range files {
		if f.IsAudio {
			infos = append(infos, GetAudioInfo(f.Path))
		}
	}

	audioJson, err := json.MarshalIndent(infos, "", "  ")
	if err != nil {
		log.Println(err)
	}

	err = tmpl.Lookup("main").Execute(w, msi{
		"fullPath":  abs,
		"path":      rel,
		"pathSplit": paths,
		"files":     files,
		"numFiles":  len(files),
		"audioJson": audioJson,
	})
	if err != nil {
		log.Println(err)
	}
}

// like join but always starts with /
func urljoin(ss ...string) string {
	return "/" + strings.TrimLeft(join(ss...), "/")
}

func join(ss ...string) string {
	clean := make([]string, 0, len(ss))
	for _, s := range ss {
		if s != "/" && len(s) > 0 {
			clean = append(clean, s)
		}
	}
	return strings.Join(clean, "/")
}

func pathSplit(path string) DirList {
	d := DirList{DirEntry{Name: ".", Url: "/", Path: "./"}}
	if path == "." {
		return d
	}

	spl := strings.Split(path, "/")

	for _, s := range spl {
		if len(s) == 0 {
			continue
		}
		prev := d[len(d)-1]
		e := DirEntry{Name: s, Url: urljoin(prev.Url, s), Path: join(prev.Path, s)}
		d = append(d, e)
	}
	return d
}

// return the rel and abs path correspondign to a url
func paths(uri string) (rel, abs string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", err
	}
	rel = filepath.Join(".", u.Path)
	abs, err = filepath.Abs(rel)
	if err != nil {
		return "", "", err
	}

	return rel, abs, nil
}

func logErr(w io.Writer, err error) {
	fmt.Fprint(w, err)
	log.Println(err)
}
