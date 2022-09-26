// Simple HTTP Form File Uploader
// Copyright (c) 2022 by Google LLC
// Fo upload file:
// $ curl -F 'file=@myfile' http://...
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	addr = flag.String("addr", ":8080", "bind address and port")
	dir  = flag.String("dir", "/tmp", "directory where upload will land")
)

func msg(w io.Writer, out string, err error) {
	m := fmt.Sprintf("%v: %v\n", out, err)
	fmt.Fprint(w, m)
	log.Print(m)
}

func main() {
	flag.Parse()
	log.Printf("Starting http listener: %v", *addr)
	if err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		w.Header().Set("Content-Type", "text/plain")
		f, h, err := r.FormFile("file")
		if err != nil {
			msg(w, "form file error", err)
			return
		}
		defer f.Close()
		o, err := os.OpenFile(filepath.Join(*dir, filepath.Base(h.Filename)), os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			msg(w, "unable to open destination file", err)
			return
		}
		defer o.Close()
		oS, err := io.Copy(o, f)
		if err != nil {
			msg(w, "unable to open destination file", err)
			return
		}
		msg(w, fmt.Sprintf("received file: %v size: %v remote: %v", o.Name(), oS, r.RemoteAddr), fmt.Errorf("ok"))
	})); err != nil {
		log.Fatalf("unable to start http server: %v", err)
	}
}
