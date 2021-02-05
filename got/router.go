package got

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sort"
)

var formStr = `<form action="/add" method="post">
KEY: <input type="text" name="key"><br/>
URL: <input type="text" name="url"><br/>
<input type="submit" value="提交">
</form>`

var store *Store

func init() {
	store = NewStore()
}

func Add(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	if url == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, formStr)
	} else {
		key := r.PostFormValue("key")
		err := store.Add(&key, &url)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, key)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	var key = r.URL.Path[1:]
	var path string
	err := store.Get(&key, &path)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	log.Println("path", path)
	http.Redirect(w, r, path, http.StatusFound)
}

func Show(w http.ResponseWriter, r *http.Request) {
	urls := store.GetUrls()
	var keys = make([]string, len(urls))
	idx := 0
	for k, _ := range urls {
		keys[idx] = k
		idx += 1
	}
	sort.Strings(keys)
	var buffer = bytes.NewBufferString("")
	for i := 0; i < len(keys); i++ {
		buffer.WriteString(keys[i])
		buffer.WriteString(" : ")
		buffer.WriteString(urls[keys[i]])
		buffer.WriteString("\n")
	}
	w.Write(buffer.Bytes())
}
