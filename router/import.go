// Overview Create a way to Declare a page which registers the
// * assets to serve
// * html head tags
package router

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Asset struct {
	Path string `json:"path"`
	Url  string `json:"url"`
	Typ  string `json:"typ"`
}

type ComponentAsset struct {
	Path     string   `json:"path"`
	IsPage   bool     `json:"isPage"`
	Children []string `json:"children"`
	Assets   []Asset  `json:"assets"`
}

// Removes . and .. from path, replaces with /
func GetUrlFromPath(path string) string {
	url := path

	if len(url) <= 1 {
		if url == "." {
			url = "/"
		}

		return url
	}

	if url[0] == '.' {
		if url[1] == '/' {
			url = url[1:]
		} else if url[1] == '.' {
			tmp := GetUrlFromPath(url[1:])
			url = tmp
		} else {
			url = url[1:]
		}
	}

	return url
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

func getEnvVarOrDefault(envVarName string, defaultVal string) string {
	v := os.Getenv(envVarName)

	if v == "" {
		v = defaultVal
	}

	return v
}

var ComponentsPath = "components/"
var PagesPath = "pages/"
var AssetsPath = "assets/"
var SupportedAssetTypes = []string{"css", "js", "scss", "png", "jpg", "jpeg", "svg"}
var TemplateFileType = "templ"
var AssetMapFileName = getEnvVarOrDefault("ASSET_MAP_FILENAME", "asset_map.json")

type ComponentMap = map[string]*ComponentAsset
type AssetMap = map[string]*Asset

func ParsePageContents(path string) (*ComponentAsset, error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	path = GetUrlFromPath(path)

	scanner := bufio.NewScanner(file)

	children := make([]string, 0)
	assets := make([]Asset, 0)

	for scanner.Scan() {
		line := scanner.Text()

		// Handle Imports
		//TODO: Cater for multi line import syntax
		if strings.Contains(line, "import") {
			if strings.Contains(line, "css") || strings.Contains(line, "js") {
				styleSheetPath := strings.Split(line, " ")[1]
				assets = append(assets, Asset{Path: styleSheetPath, Typ: "css", Url: GetUrlFromPath(styleSheetPath)})
			} else if strings.Contains(line, "components") {
				childPath := strings.Split(line, " ")[1]
				componentIndex := strings.Index(childPath, ComponentsPath)

				if componentIndex == -1 {
					return nil, fmt.Errorf("invalid path for component %s %s", childPath, line)
				}

				if childPath[0] == '.' {
					childPath = childPath[1:]
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
				Url:  GetUrlFromPath(path),
				Typ:  fileType,
			}

			assets = append(assets, asset)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &ComponentAsset{
		Path:     path,
		Assets:   assets,
		IsPage:   strings.Contains(path, "pages"),
		Children: children,
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

		if fullPath[0] == '.' && (len(fullPath) > 1 && fullPath[1] != '/' && fullPath[1] != '.') {
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

			(*compMap)[GetUrlFromPath(fullPath[:index-1])] = compAsset
		} else if slices.Contains(SupportedAssetTypes, fileType) {
			asset := Asset{
				Path: fullPath,
				Typ:  fileType,
				Url:  GetUrlFromPath(fullPath),
			}

			(*assetMap)[fullPath] = &asset
		} else {
			return fmt.Errorf("error handling unknown file: \nname: %s\npath: %s", fileName, path)
		}
	}

	return nil
}

func GetChildAssets(compMap *ComponentMap, childPath string, assetMap *AssetMap) error {
	ok := false
	child := ComponentAsset{}

	if childPath[0] == '.' {
		childPath = childPath[1:]
	}

	for path, compAsset := range *compMap {
		splitPath := strings.Split(path, "/")

		// Check if final word of path is templ
		// if strings.Contains(splitPath[len(splitPath)-1], "templ") {
		// 	index := strings.Index(path, splitPath[len(splitPath)-1])
		// 	path = path[:index-1]
		// }
		if len(path) > 5 && path[len(path)-5:] == "templ" {
			index := strings.Index(path, splitPath[len(splitPath)-1])
			if index == -1 {
				return fmt.Errorf("error getting child assets\npath: %s\ncompMap: %v", path, *compMap)
			}
			path = path[:index]
		}

		splitPath = strings.Split(childPath, ".")

		if len(path) > 5 && path[len(path)-5:] == "templ" {
			index := strings.Index(path, splitPath[len(splitPath)-1])
			if index == -1 {
				return fmt.Errorf("error getting child assets\npath: %s\ncompMap: %v", childPath, *compMap)
			}
			childPath = childPath[:index-1]
		}

		if strings.Contains(path, childPath) || strings.Contains(childPath, path) {
			child = *compAsset
			ok = true
		}
	}

	if !ok {
		return fmt.Errorf("child does not exist in component map %s", childPath)
	}

	for _, nestedChildPath := range child.Children {
		err := GetChildAssets(compMap, nestedChildPath, assetMap)

		if err != nil {
			return err
		}
	}

	for _, asset := range child.Assets {
		if _, ok := (*assetMap)[asset.Path]; !ok {
			(*assetMap)[asset.Path] = &asset
		}

	}

	return nil
}

// Create a function to read map and create head for page
func CreatePageHead(compMap *ComponentMap, path string) (AssetMap, error) {
	compAsset := (*compMap)[path]

	if !compAsset.IsPage {
		return nil, fmt.Errorf("component at path %s is not a page\ncomp: %v", path, *compAsset)
	}

	assetMap := make(AssetMap)

	for _, childPath := range compAsset.Children {
		err := GetChildAssets(compMap, childPath, &assetMap)

		if err != nil {
			return nil, err
		}
	}

	for _, asset := range compAsset.Assets {
		assetMap[asset.Path] = &asset
	}

	return assetMap, nil
}

// Uses compononents, pages and assets path to load required imports
func LoadImports(rootDir string) (ComponentMap, AssetMap) {
	compMap := make(ComponentMap)
	assetMap := make(AssetMap)

	err := RegisterAssets(rootDir, true, &compMap, &assetMap)

	if err != nil {
		log.Fatal("error registering assets: ", err.Error())
	}

	for _, asset := range assetMap {
		assetIndex := strings.Index(asset.Path, AssetsPath)
		componentIndex := strings.Index(asset.Path, ComponentsPath)
		pageIndex := strings.Index(asset.Path, PagesPath)

		assetUrl := ""

		if assetIndex != -1 {
			assetUrl = asset.Path[assetIndex:]
		}

		if componentIndex != -1 {
			assetUrl = asset.Path[componentIndex+len(ComponentsPath):]
		}

		if pageIndex != -1 {
			assetUrl = asset.Path[pageIndex+len(PagesPath):]
		}

		if assetUrl == "" {
			log.Fatal("asset url not defined for asset ", *asset)
		}

		asset.Url = assetUrl
	}

	return compMap, assetMap
}

func CreateAssetsFile(path string) {
	compMap, assetMap := LoadImports(path)

	file, err := os.Create(AssetMapFileName)

	if err != nil {
		log.Fatalf("error creating file: %s", err.Error())
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(AssetFile{
		CompMap:  compMap,
		AssetMap: assetMap,
	})

	if err != nil {
		log.Fatalf("error encoding file contents: %s", err.Error())
	}
}

type AssetFile struct {
	CompMap  ComponentMap `json:"compMap"`
	AssetMap AssetMap     `json:"assetMap"`
}

func ReadAssetsFile(assetFilePath string) (ComponentMap, AssetMap) {
	file, err := os.Open(assetFilePath)

	if err != nil {
		log.Fatal("Error opening file:", err.Error())
	}
	defer file.Close()

	assetFile := AssetFile{}

	err = json.NewDecoder(file).Decode(&assetFile)

	if err != nil {
		log.Fatal("Error decoding JSON:", err.Error())
	}

	return assetFile.CompMap, assetFile.AssetMap
}
