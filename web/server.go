package main

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"github.com/kardianos/service"
	"github.com/mlctrez/servicego"
	"github.com/mlctrez/wasmexec"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Service struct {
	servicego.Defaults
	httpServer *http.Server
}

func (s *Service) Start(_ service.Service) error {

	httpServer := &http.Server{}
	httpServer.Addr = os.Getenv("ADDRESS")
	if httpServer.Addr == "" {
		httpServer.Addr = ":8080"
	}

	mux := http.NewServeMux()
	httpServer.Handler = mux

	files := http.FileServer(http.FS(fs))
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/favicon.ico":
			w.WriteHeader(404)
		case "/app.js":
			writeAppJs(w)
		default:
			files.ServeHTTP(w, r)
		}
	}
	mux.HandleFunc("/", handler)
	s.httpServer = httpServer

	go func() {
		_ = httpServer.ListenAndServe()
	}()
	return nil
}

func (s *Service) Stop(_ service.Service) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(context.Background())
	}
	return nil
}

func main() {
	servicego.Run(&Service{})
}

//go:embed *.js *.html *.css *.wasm
var fs embed.FS

func writeAppJs(w http.ResponseWriter) {

	var err error
	var content []byte

	content, err = wasmexec.Current()
	if err != nil {
		w.WriteHeader(500)
		log.Printf("error getting current wasmexec: %v", err)
		return
	}

	var launch []byte
	launch, err = fs.ReadFile("app_launch.js")
	if err != nil {
		w.WriteHeader(500)
		log.Printf("error reading app_launch.js: %v", err)
		return
	}

	_, err = w.Write(append(append(content, []byte("\n\n")...), launch...))
	if err != nil {
		log.Printf("write error: %v", err)
		return
	}

}

func buildWasm() (err error) {
	command := exec.Command("go", "build", "-o", "web/app.wasm", ".")
	command.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	var output []byte
	if output, err = command.CombinedOutput(); err != nil {
		return fmt.Errorf("%s \n\n %s", string(output), err.Error())
	}
	var stat os.FileInfo
	if stat, err = os.Stat("web/app.wasm"); err != nil {
		return err
	}
	fmt.Println(stat.Size())
	return nil
}
