// Overview Create a way to Declare a page which registers the
// * assets to serve
// * html head tags
package router

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Asset struct {
	Path string
	Typ  string
}

type ComponentAsset struct {
	path     string
	isPage   bool
	children []string
	assets   []Asset
}

// Combine 2 url paths, removing the trailing / and catering for overlapping /s
// for example
// /base/ + /append/ = /base/append
func AppendPath(basePath string, path string) string {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if basePath[len(basePath)-1] == '/' {
		return basePath + path
	}

	return basePath + "/" + path
}

var componentsPath = "components/"
var pagesPath = "pages/"
var assetsPath = "assets/"
var SupportedAssetTypes = []string{"css", "js", "scss", "png", "jpg", "jpeg", "svg"}
var TemplateFileType = "templ"

type ComponentMap = map[string]ComponentAsset
type AssetMap = map[string]Asset

func ParsePageContents(path string) (*ComponentAsset, error) {

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	children := make([]string, 0)
	assets := make([]Asset, 0)

	for scanner.Scan() {
		line := scanner.Text()

		// Handle Imports
		//TODO: Cater for multi line import syntax
		if strings.Contains(line, "import") {
			if strings.Contains(line, "css") {
				styleSheetPath := strings.Split(line, " ")[1]
				assets = append(assets, Asset{Path: styleSheetPath, Typ: "css"})
			} else if strings.Contains(line, "components") {
				childPath := strings.Split(line, " ")[1]
				componentIndex := strings.Index(childPath, componentsPath)

				if componentIndex == -1 {
					return nil, fmt.Errorf("invalid path for component %s %s", childPath, line)
				}
				children = append(children, childPath)
			}

			continue
		}

		// Check for All Assets inline, like image tags
		// TODO: Need to cater for certain things being relative paths
		if strings.Contains(line, "/assets/") {
			splitLine := strings.Split(line, "\"")
			path := ""

			for _, str := range splitLine {
				if strings.Contains(str, "/assets/") {
					path = str
				}
			}

			if path == "" {
				return nil, fmt.Errorf("issue finding relevant asset path")
			}

			splitFileStr := strings.Split(path, "/")
			fileName := splitFileStr[len(splitFileStr)-1]
			fileType := strings.Split(fileName, ".")[1]

			asset := Asset{
				Path: path,
				Typ:  fileType,
			}

			assets = append(assets, asset)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &ComponentAsset{
		path:     path,
		assets:   assets,
		isPage:   strings.Contains(path, "pages"),
		children: children,
	}, nil
}

// Reads through all pages + components and creates map
// Read imports for components used
// Read through components for assets used
func RegisterAssets(path string, recursive bool, compMap *ComponentMap, assetMap *AssetMap) error {
	dir, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	for _, file := range dir {

		fileName := file.Name()
		fullPath := AppendPath(path, fileName)

		if file.IsDir() {
			if !recursive {
				continue
			}

			err = RegisterAssets(fullPath, true, compMap, assetMap)

			if err != nil {
				return err
			}

			continue
		}

		splitFileStr := strings.Split(fileName, ".")

		if len(splitFileStr) < 2 {
			continue
		}

		fileType := splitFileStr[1]

		if !slices.Contains(SupportedAssetTypes, fileType) && fileType != "templ" {
			continue
		}

		// TODO: Note, certain files have . at the front
		// TODO: Need to remove github.com from component path

		if fullPath[0] == '.' {
			fullPath = fullPath[1:]
		}

		if fullPath[0] == '/' && len(fullPath) > 1 {
			fullPath = fullPath[1:]
		}

		_, exists := (*compMap)[fullPath]

		if exists {
			continue
		}

		_, exists = (*assetMap)[fullPath]

		if exists {
			continue
		}

		if fileType == "templ" {
			compAsset, err := ParsePageContents(fullPath)

			if err != nil {
				return err
			}

			index := strings.Index(fullPath, fileName)

			(*compMap)[fullPath[:index-1]] = *compAsset
		} else if slices.Contains(SupportedAssetTypes, fileType) {

			(*assetMap)[fullPath] = Asset{
				Path: fullPath,
				Typ:  fileType,
			}
		}
	}

	return nil
}

func GetChildAssets(compMap *ComponentMap, childPath string, assetMap *AssetMap) error {
	ok := false
	child := ComponentAsset{}

	for path, compAsset := range *compMap {
		splitPath := strings.Split(path, "/")

		if strings.Contains(splitPath[len(splitPath)-1], "templ") {
			index := strings.Index(path, splitPath[len(splitPath)-1])
			path = path[:index-1]
		}

		splitPath = strings.Split(childPath, ".")

		if strings.Contains(splitPath[len(splitPath)-1], "templ") {
			index := strings.Index(path, splitPath[len(splitPath)-1])
			childPath = childPath[:index-1]
		}

		if strings.Contains(path, childPath) || strings.Contains(childPath, path) {
			child = compAsset
			ok = true
		}
	}

	if !ok {
		fmt.Println("error map", childPath, compMap)
		return fmt.Errorf("child does not exist in component map %s", childPath)
	}

	for _, nestedChildPath := range child.children {
		err := GetChildAssets(compMap, nestedChildPath, assetMap)

		if err != nil {
			return err
		}
	}

	for _, asset := range child.assets {
		if _, ok := (*assetMap)[asset.Path]; !ok {
			(*assetMap)[asset.Path] = asset
		}

	}

	return nil
}

// Create a function to read map and create head for page
func CreatePageHead(compMap *ComponentMap, path string) (AssetMap, error) {
	compAsset := (*compMap)[path]

	if !compAsset.isPage {
		fmt.Println("failing map", compAsset, "path is ", path)
		return nil, fmt.Errorf("component at path %s is not a page", path)
	}

	assetMap := make(AssetMap)

	for _, childPath := range compAsset.children {
		err := GetChildAssets(compMap, childPath, &assetMap)

		if err != nil {
			return nil, err
		}
	}

	for _, asset := range compAsset.assets {
		assetMap[asset.Path] = asset
	}

	return assetMap, nil
}

// Uses compononents, pages and assets path to load required imports
func LoadImports(rootDir string, r Router) ComponentMap {
	compMap := make(ComponentMap)
	assetMap := make(AssetMap)

	err := RegisterAssets(rootDir, true, &compMap, &assetMap)

	if err != nil {
		os.Exit(1)
	}

	for _, asset := range assetMap {
		assetIndex := strings.Index(asset.Path, assetsPath)
		componentIndex := strings.Index(asset.Path, componentsPath)
		pageIndex := strings.Index(asset.Path, pagesPath)

		assetUrl := ""

		if assetIndex != -1 {
			assetUrl = asset.Path[assetIndex:]
		}

		if componentIndex != -1 {
			assetUrl = asset.Path[componentIndex+len(componentsPath):]
		}

		if pageIndex != -1 {
			assetUrl = asset.Path[pageIndex+len(pagesPath):]
		}

		if assetUrl == "" {
			os.Exit(1)
		}

		r.Serve(assetUrl, asset.Path, &RouteOptions{})
	}

	return compMap
}
