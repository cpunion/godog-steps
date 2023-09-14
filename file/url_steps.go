package file

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ctxServerKey struct {
}

type ctxServerState struct {
	server  *http.Server
	mapping map[string]string
}

func getServerState(ctx context.Context) *ctxServerState {
	svr := ctx.Value(ctxServerKey{})
	return svr.(*ctxServerState)
}

func getServer(ctx context.Context) *http.Server {
	svr := getServerState(ctx)
	return svr.server
}

type mappingFS struct {
	http.FileSystem
	mappings map[string]string
}

func NewMappingFS(fs http.FileSystem, mappings map[string]string) http.FileSystem {
	return &mappingFS{FileSystem: fs, mappings: mappings}
}

func (fs *mappingFS) Open(name string) (http.File, error) {
	if mapped, ok := fs.mappings[name]; ok {
		name = mapped
	}
	return fs.FileSystem.Open(name)
}

func iHaveAUrlFromFile(ctx context.Context, link, file string) (context.Context, error) {
	// Mock url with a local file by http://localhost:port
	u, err := url.Parse(link)
	if err != nil {
		return ctx, errors.New("invalid url")
	}
	isLocalhost := strings.HasPrefix(u.Host, "localhost:") || strings.HasPrefix(u.Host, "127.0.0.1:")
	if u.Scheme != "http" || !isLocalhost {
		return ctx, errors.New("invalid url, must starts with http://localhost:port or http://127.0.0.1:port")
	}
	port, err := strconv.Atoi(strings.Split(u.Host, ":")[1])
	if err != nil {
		return ctx, err
	}
	handler := http.NewServeMux()
	mappings := map[string]string{u.Path: file}
	handler.Handle(u.Path, http.FileServer(NewMappingFS(http.FS(assetsFS), mappings)))
	addr := ":" + strconv.Itoa(port)
	server := &http.Server{Addr: addr, Handler: handler}
	getServerState(ctx).server = server
	go func() {
		server.ListenAndServe()
	}()
	return ctx, nil
}
