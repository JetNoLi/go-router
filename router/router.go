package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/jetnoli/go-router/utils"
)

// TODO: REMOVE AND REPLACE WITH BUILT IN FUNC FROM HTTP LIB
// TODO: Move types to their own files
// TODO: Create custom Interface for Req, Res (ctx ?)
// TODO: Add Request ID
// TODO: Add in api version flow to router?
// TODO: Add in doc comments

func ReadData(path string) (data []byte, err error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return data, err
	}

	return os.ReadFile(absPath)
}

type RouterOptions struct {
	ExactPathsOnly        bool // Appends the {$} for all paths in router
	PreHandlerMiddleware  []MiddlewareHandler
	PostHandlerMiddleware []MiddlewareHandler
}

type RouteOptions struct {
	PreHandlerMiddleware  []MiddlewareHandler
	PostHandlerMiddleware []MiddlewareHandler
}

type MiddlewareHandler = func(w *http.ResponseWriter, r *http.Request)

type Router struct {
	Path    string
	Mux     *http.ServeMux
	Options RouterOptions
}

// TODO: Add Global Response Headers
func CreateRouter(path string, options RouterOptions) *Router {
	router := &Router{}

	router.Mux = http.NewServeMux()

	router.Path = path
	router.Options = options

	if options.PostHandlerMiddleware == nil {
		router.Options.PostHandlerMiddleware = make([]MiddlewareHandler, 0)
	}

	if options.PreHandlerMiddleware == nil {
		router.Options.PreHandlerMiddleware = make([]MiddlewareHandler, 0)
	}

	return router
}

func (router Router) CreateRoute(path string, method string) string {
	pathEndString := ""

	if router.Options.ExactPathsOnly {
		end := len(path) - 1

		if len(path) > 3 && path[end-3:end] == "{$}" {
			pathEndString = ""
		} else if path[len(path)-1] == '/' {
			pathEndString = "{$}"
		} else {
			pathEndString = "/{$}"
		}
	}

	url := method + " " + router.Path

	// To avoid double // in request e.g. GET //path-name
	if router.Path[len(router.Path)-1] == '/' && path[0] == '/' {
		url += path[1:]
	} else {
		url += path
	}

	fmt.Println("Registering: " + url + pathEndString)

	return url + pathEndString
}

// TODO: Add config to only run these in DEBUG mode
func (router Router) ExecuteWithMiddleware(w *http.ResponseWriter, r *http.Request, handler http.HandlerFunc, routeOptions *RouteOptions) {

	//TODO: Make sure doesn't pass by reference
	preHandlerMiddleware := router.Options.PreHandlerMiddleware[:]
	postHandlerMiddleware := router.Options.PostHandlerMiddleware[:]

	if routeOptions != nil {
		if routeOptions.PreHandlerMiddleware != nil {
			preHandlerMiddleware = append(preHandlerMiddleware, routeOptions.PreHandlerMiddleware...)

		}

		if routeOptions.PostHandlerMiddleware != nil {
			postHandlerMiddleware = append(postHandlerMiddleware, routeOptions.PostHandlerMiddleware...)
		}

	}

	for _, middleware := range preHandlerMiddleware {
		middleware(w, r)

		if r.Context().Err() != nil {
			return
		}
	}

	handlerName := utils.GetFunctionName(handler)
	fmt.Println("executing handler ", handlerName)

	handler(*w, r)

	for _, middleware := range postHandlerMiddleware {
		if r.Context().Err() != nil {
			return
		}

		fmt.Printf("middleware applied %s", utils.GetFunctionName(middleware))
		middleware(w, r)
	}

}

func (router Router) HandleFunc(path string, handler http.HandlerFunc, options *RouteOptions) {
	router.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		rCtxCopy, cancel := context.WithCancel(r.Context())
		*r = *r.WithContext(context.WithValue(rCtxCopy, utils.CancelRequestKey, cancel))

		fmt.Println("Serving path ", path)
		router.ExecuteWithMiddleware(&w, r, handler, options)

	})
}

func (router Router) Handle(path string, mux *http.ServeMux) {
	router.Mux.Handle(path, &RouteHandler{
		ChildMux: mux,
		Router:   &router,
	})
}

func (router Router) Get(path string, handler http.HandlerFunc, options *RouteOptions) {
	route := router.CreateRoute(path, "GET")
	router.HandleFunc(route, handler, options)
}

func (router Router) Post(path string, handler http.HandlerFunc, options *RouteOptions) {
	route := router.CreateRoute(path, "POST")
	router.HandleFunc(route, handler, options)
}

func (router Router) Delete(path string, handler http.HandlerFunc, options *RouteOptions) {
	route := router.CreateRoute(path, "DELETE")
	router.HandleFunc(route, handler, options)
}

func (router Router) Put(path string, handler http.HandlerFunc, options *RouteOptions) {
	route := router.CreateRoute(path, "PUT")
	router.HandleFunc(route, handler, options)
}

func (router Router) Patch(path string, handler http.HandlerFunc, options *RouteOptions) {
	route := router.CreateRoute(path, "PATCH")
	router.HandleFunc(route, handler, options)
}

// Templating

// Serve file at the given filepath relative to app
func (router Router) Serve(path string, filePath string, options *RouteOptions) {
	route := router.CreateRoute(path, "GET")

	router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		data, err := ReadData(filePath)

		if err != nil {
			http.Error(w, "Error Reading file:\n"+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", http.DetectContentType(data))

		_, err = w.Write(data)

		if err != nil {
			http.Error(w, fmt.Sprintf("error serving file: %s", filePath), http.StatusInternalServerError)
		}

	}, options)
}

type ServeDirOptions struct {
	IncludedExtensions         []string
	RoutePathContainsExtension bool
	Recursive                  bool
}

// Serves all html in given directory relative to app
// Include trailing slash in dir name
// TODO: rewrite with built in http functionality
func (router Router) ServeDir(baseUrlPath string, dirPath string, options *ServeDirOptions) {
	absPath, err := filepath.Abs(dirPath)

	if err != nil {
		panic("error reading dir: " + err.Error())
	}

	files, err := os.ReadDir(absPath)

	if err != nil {
		panic("failed to read dir for html pages: " + err.Error())
	}

	for _, file := range files {

		if file.IsDir() {
			if options.Recursive {
				router.ServeDir(baseUrlPath+file.Name()+"/", dirPath+file.Name()+"/", options)
			}

			continue
		}

		fileName := file.Name()
		fileExtention := filepath.Ext(fileName)

		if len(options.IncludedExtensions) > 0 && !slices.Contains(options.IncludedExtensions, fileExtention) {
			continue
		}

		filePath := absPath + "/" + fileName

		route := baseUrlPath + fileName

		if !options.RoutePathContainsExtension {
			route = route[:len(route)-len(fileExtention)]
		}

		router.Serve(route, filePath, &RouteOptions{})
	}
}
