# go-net-utils
A small collection of utilites and wrappers to make building http applications with go simpler


## Overview
All go router functionality will work in combination with the standard net/http lib and thus is compatible with most popular HTTP Frameworks.

The API is inspired by express in structure and allows for easy separation of endpoints via Routers, with a simple interface for handling requests


## Basic Server Example
```golang

// Define Health Check Handler
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Hello World!`))
}

// Create Base Router to be Used as Server
router := Router.CreateRouter("/", Router.RouterOptions{
    // Attach Middleware if Required
	PreHandlerMiddleware: []Router.MiddlewareHandler{middleware.DecodeToken},
})

// Create a Get Request Handler at /health path
router.Get("/health", HealthCheck, &Router.RouteOptions{})

// Define Server with Standard Http Library
server := http.Server{
	Addr:         ":" + config.Port,
	ReadTimeout:  60 * time.Second,
	WriteTimeout: 60 * time.Second,
	Handler:      router.Mux,
}

fmt.Println("Starting Server on http://localhost:" + config.Port)

// Start Server
log.Fatal(server.ListenAndServe())

```

## Basic Router Example
```golang

// Define Router
func AuthRouter() *http.ServeMux {
	router := Router.CreateRouter("/auth", Router.RouterOptions{})

	router.Post("/login", handlers.SignInHtmx, &Router.RouteOptions{})
	router.Post("/signup", handlers.SignUpHtmx, &Router.RouteOptions{})

	return router.Mux
}

// Attach Router to Server
router.Handle("/auth/", AuthRouter())    

```

## Templ Serve Example
```golang
func HTMLRouter() *http.ServeMux {
    router := Router.CreateRouter("/", Router.RouterOptions{
        ExactPathsOnly: true,
    })

    // Serve Styles for Pages
    router.ServeDir("/", "view/pages/", &Router.ServeDirOptions{
        IncludedExtensions:         []string{".css"},
        Recursive:                  true,
        RoutePathContainsExtension: true,
    })

    // Serve Styles for Components
    router.ServeDir("/", "view/components/", &Router.ServeDirOptions{
        IncludedExtensions:         []string{".css"},
        Recursive:                  true,
        RoutePathContainsExtension: true,
    })

    // Serve Assets for Components
    router.ServeDir("/assets/", "assets/", &Router.ServeDirOptions{
        Recursive:                  true,
        RoutePathContainsExtension: true,
    })

    router.ServeTempl(map[string]*Router.TemplPage{
        "/": {
            PageComponent: home.Page(),
            Options: &Router.RouteOptions{
                PreHandlerMiddleware: []Router.MiddlewareHandler{middleware.CheckAuthorization},
            },
        },
        "/login": {
            PageComponent: login.Page(),
        },
        "/signup": {
            PageComponent: signup.Page(),
        },
        "/admin/": {
            PageComponent: admin.Page(admin.PageUrlParams{}),
        },
    })

    return router.Mux
}


```