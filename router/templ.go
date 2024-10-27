package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

// TODO: Pass component map in ctx
// TODO: Allow a way to pass in any data the component may need?
// Mapybe should be a ServeTemplHandler Option?
func (r Router) ServeTempl(route string, comp templ.Component, compMap *ComponentMap) error {

	// assumption: route contains page name
	routeSplit := strings.Split(route, "/")
	localPathName := routeSplit[len(routeSplit)-1]
	routeSplitLen := len(routeSplit)

	// check if route param
	if route == "/" {
		localPathName = "home"
	}

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

	localPath := ""

	for compPath, compAsset := range *compMap {
		if !compAsset.isPage {
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

	fmt.Println("map a", assetMap)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err = Render(w, r, comp, assetMap)

		if err != nil {
			http.Error(w, "error serving page "+err.Error(), http.StatusInternalServerError)
		}

	}, &RouteOptions{})

	return nil
}
