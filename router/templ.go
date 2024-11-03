package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

type ServeTemplOptions struct {
	HeadData *templ.Component
}

type ServeTemplOptsFunc = func(opts *ServeTemplOptions)

func WithHead(headData *templ.Component) ServeTemplOptsFunc {
	return func(opts *ServeTemplOptions) {
		if opts == nil {
			opts = &ServeTemplOptions{}
		}

		opts.HeadData = headData
	}
}

func (r Router) ServeTempl(route string, comp templ.Component, compMap *ComponentMap, optsFuncs ...ServeTemplOptsFunc) error {

	// assumption: route contains page name
	routeSplit := strings.Split(route, "/")
	localPathName := routeSplit[len(routeSplit)-1]
	routeSplitLen := len(routeSplit)

	// default / path to home component
	// TODO: allow for override
	if route == "/" {
		localPathName = "home"
	}

	// check for route params, defined by {}, e.g. /user/{id}
	if localPathName != "home" {
		for i := 1; i < routeSplitLen; i++ {
			path := routeSplit[routeSplitLen-i]

			if !strings.Contains(path, "{") {
				localPathName = path
				break
			}

			if i == routeSplitLen-1 {
				return fmt.Errorf("no matching component path found for %s %v", route, compMap)
			}
		}
	}

	// TODO: find a better flow here
	localPath := ""

	for compPath, compAsset := range *compMap {
		if !compAsset.IsPage {
			continue
		}

		// TODO: Need to consider how overlapping names can affect this contains check
		if strings.Contains(compPath, localPathName) {
			localPath = compPath
		}
	}

	if localPath == "" {
		return fmt.Errorf("no local path found for local path name %s on route %s", localPathName, route)
	}

	assetMap, err := CreatePageHead(compMap, localPath)

	if err != nil {
		return fmt.Errorf("error creating page header content from asset map " + err.Error())
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err = Render(w, r, comp, assetMap)

		if err != nil {
			http.Error(w, "error serving page "+err.Error(), http.StatusInternalServerError)
		}

	}, &RouteOptions{})

	return nil
}
