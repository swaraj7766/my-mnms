package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed build
var staticFS embed.FS
var Version string

func main() {
	flagversion := flag.Bool("version", false, "print version")
	addr := flag.String("addr", ":9000", "address to bind and serve")
	flag.Parse()
	if *flagversion {
		info, _ := debug.ReadBuildInfo()
		fmt.Fprintln(os.Stderr, Version)
		fmt.Fprintln(os.Stderr, info.GoVersion, info.Main, info.Settings)
		os.Exit(0)
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	AddRoutes(router)
	err := router.Run(*addr)
	if err != nil {
		fmt.Print(err)
	}
}

// AddRoutes serves the static file system for the UI React App.
func AddRoutes(router gin.IRouter) {
	embeddedBuildFolder := newStaticFileSystem()
	fallbackFileSystem := newFallbackFileSystem(embeddedBuildFolder)
	router.Use(static.Serve("/", embeddedBuildFolder))
	router.Use(static.Serve("/", fallbackFileSystem))
}

// ----------------------------------------------------------------------
// staticFileSystem serves files out of the embedded build folder

type staticFileSystem struct {
	http.FileSystem
}

var _ static.ServeFileSystem = (*staticFileSystem)(nil)

func newStaticFileSystem() *staticFileSystem {
	sub, err := fs.Sub(staticFS, "build")
	if err != nil {
		panic(err)
	}

	return &staticFileSystem{
		FileSystem: http.FS(sub),
	}
}

func (s *staticFileSystem) Exists(prefix string, path string) bool {
	buildpath := fmt.Sprintf("build%s", path)

	// support for folders
	if strings.HasSuffix(path, "/") {
		_, err := staticFS.ReadDir(strings.TrimSuffix(buildpath, "/"))
		return err == nil
	}

	// support for files
	f, err := staticFS.Open(buildpath)
	if f != nil {
		_ = f.Close()
	}
	return err == nil
}

// ----------------------------------------------------------------------
// fallbackFileSystem wraps a staticFileSystem and always serves /index.html
type fallbackFileSystem struct {
	staticFileSystem *staticFileSystem
}

var (
	_ static.ServeFileSystem = (*fallbackFileSystem)(nil)
	_ http.FileSystem        = (*fallbackFileSystem)(nil)
)

func newFallbackFileSystem(staticFileSystem *staticFileSystem) *fallbackFileSystem {
	return &fallbackFileSystem{
		staticFileSystem: staticFileSystem,
	}
}

func (f *fallbackFileSystem) Open(path string) (http.File, error) {
	return f.staticFileSystem.Open("/index.html")
}

func (f *fallbackFileSystem) Exists(prefix string, path string) bool {
	return true
}
