package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"flag"
	"time"
	"os"
	"path/filepath"
	"strings"
trie	"github.com/henesy/fstrie"
)

// Website root -- could easily extend this and the `*file`s to reference the subdomain root used
var www string


// Virtual files do ??, physical files print their bytes to string
type stringable interface {
	Read(string) []byte
}


// Represents a virtual file
type virtFile struct {
}

// Writes the default virtual file's bytes -- Replace to override
func (f *virtFile) Read(request string) []byte {
	return []byte(virtString)
}

// Make a new default virtFile
func mkVF() *virtFile {
	f := new(virtFile)
	//f.Read = defaultRead
	return f
}


// Represents a real file
type realFile struct {
	path	string
}

// Writes the real file's bytes
func (f *realFile) Read(request string) []byte {
	content, err := ioutil.ReadFile(www + f.path)
	if err != nil {
		log.Print("Error, can't read file: ", err)
		return []byte("<html><p1>ERROR</p1></br>" + request + "</html>")
	}
	return content
}

// Make a new default realFile
func mkRF(p string) *realFile {
	f := new(realFile)
	f.path = p
	return f
}


func walkRoot(path string) ([]string, error) {
	fileList := make([]string, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		// Fix names -- slow
		//log.Print(strings.SplitN(path, "/", 2))
		path = strings.SplitN(path, "/", 2)[1]
		path = "/" + path
		fileList = append(fileList, path)
		return err
	})
	
	if err != nil {
		log.Print("Error walking: ", err)
		return []string{}, nil
	}

	return fileList[1:], nil
}


/* A simple web-app using a trie for its vfs */
func main() {
	flag.StringVar(&www, "r", "./www", "Set website root directory.")
	flag.Parse()

	t := trie.New()
	root := t.Find("/")
	// Make root compliant
	root.Data = mkVF()
	
	// Set up handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Maybe sanitize paths
		log.Print("Got: ", path)
		if n := t.Find(path); n != nil {
			w.Write(n.Data.(stringable).Read(path))
		} else {
			w.Write(page404)
		}
	}
	
	// Set up updater
	// Walk through directory and init new trie elements
	go func() {
		for {
			// refine
			list, _ := walkRoot(www)
			//log.Print("All files: ", list)
			for _, v := range list {
				if n :=  t.Find(v); n == nil{
					// Make the rfile
					rf := mkRF(v)
					t.Add(v, rf)
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()

	// Set up web server
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}


// String for virtual files to print by default
var virtString string = `<html>
<p2> Default virtual page. </p2>
</html>
`

// 404 page
var page404 []byte = []byte(`<!DOCTYPE html>
<html><body>
<p1>404 Error</p1></br>
<p3>OOPSIE WOOPSIE!!</p3>
</body></html>
`)
