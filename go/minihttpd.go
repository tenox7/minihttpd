package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	path     = flag.String("path", ".", "the path to the directory to serve")
	httpAddr = flag.String("addr", ":8080", "address/port to bind to")
	httpMaxf = flag.Int64("maxf", 10<<20, "max file size for http form parsing")
)

func httpUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(*httpMaxf)
	f, h, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "<HTML><BODY>\n<FORM ACTION=\"/upload\" METHOD=\"POST\" ENCTYPE=\"multipart/form-data\">\n"+
			"<INPUT TYPE=\"FILE\" NAME=\"file\"><INPUT TYPE=\"SUBMIT\" NAME=\"Upload\" VALUE=\"Upload\"></FORM>\n"+
			"<P><PRE>curl -F \"file=@yourfile.dat\" http://%v/upload</PRE></P></BODY></HTML>\n",
			r.Context().Value(http.LocalAddrContextKey).(net.Addr))
		return
	}
	defer f.Close()
	o, err := os.OpenFile(filepath.Join(*path, filepath.Base(h.Filename)), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer o.Close()
	oS, err := io.Copy(o, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "received file: %v size: %v remote: %v", o.Name(), oS, r.RemoteAddr)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<HTML><BODY>\n<FORM ACTION=\"/upload\" METHOD=\"POST\" ENCTYPE=\"multipart/form-data\">\n"+
		"<INPUT TYPE=\"FILE\" NAME=\"file\"><INPUT TYPE=\"SUBMIT\" NAME=\"Upload\" VALUE=\"Upload\"></FORM>\n"+
		"<P><PRE>curl -F \"file=@yourfile.dat\" http://%v</PRE></P>\n",
		r.Context().Value(http.LocalAddrContextKey).(net.Addr))

	dir, err := os.ReadDir(*path)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		log.Printf("Error: %v", err)
		return
	}
	for _, d := range dir {
		fmt.Fprintf(w, "<A HREF=\"%v\">%v</A><BR>\n", d.Name(), d.Name())
	}
	fmt.Fprintf(w, "</BODY></HTML>\n")
}

func startHttp() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", httpUpload)
	fmt.Printf("Serving %s on HTTP %s\n", *path, *httpAddr)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	startHttp()
}
