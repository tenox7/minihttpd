// Simple HTTP Form File Upload Handler
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

func msg(w http.ResponseWriter, out string, err error, code int) {
	m := fmt.Sprintf("%v: %v\n", out, err)
	http.Error(w, m, code)
	log.Print(m)
}

func main() {
	flag.Parse()
	log.Printf("Starting http listener: %v", *addr)
	if err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		w.Header().Set("Content-Type", "text/html")
		f, h, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintf(w, "<HTML><BODY>\n<FORM ACTION=\"/\" METHOD=\"POST\" ENCTYPE=\"multipart/form-data\">\n"+
				"<INPUT TYPE=\"FILE\" NAME=\"file\"><INPUT TYPE=\"SUBMIT\" NAME=\"Upload\" VALUE=\"Upload\">\n"+
				"</FORM></BODY></HTML>\n",
			)
			return
		}
		defer f.Close()
		o, err := os.OpenFile(filepath.Join(*dir, filepath.Base(h.Filename)), os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			msg(w, "unable to open destination file", err, http.StatusBadRequest)
			return
		}
		defer o.Close()
		oS, err := io.Copy(o, f)
		if err != nil {
			msg(w, "unable to copy destination file", err, http.StatusBadRequest)
			return
		}
		msg(w, fmt.Sprintf("received file: %v size: %v remote: %v", o.Name(), oS, r.RemoteAddr), fmt.Errorf("ok"), http.StatusOK)
	})); err != nil {
		log.Fatalf("unable to start http server: %v", err)
	}
}
