# Go Router Improvements

## Version 0.1

### Router
* All fetch.go TODOs
* All import.go TODOs
* All router.go TODOs
* All templ.go TODOs
* All utils TODOs

* Clean Up File Reading
  * Make Consistent Flow
  * Use of AbsPath vs Relative in actual builds
* Improve Env File Reading Flow
* Replace Serve and Serve Dir with Built In Functionality
* Setup Flow for DEBUG mode
* Cleanup Logs
* Review How Asset Paths are Passed, Perhaps Should be Tied to Context
* Create Server Functionality in Router
* Update ReadMe

### CLI
* Improve Dummy API UI
* Add Basic BubbleTea Interface
* Add Config to Configure HTMX as Asset in Install Flow
* Allow Root File Path to be Chosen
* Add Docker 
  * Run CLI ass a Docker App
  * Add Docker File and Docker Compose to Static so it's included in template repo 

## Version 0.2

### Router
* Add in Flow for API Versioning
* Use FOP Structure for Config
* Implement Ctx Object
* Add SSR to Serve Templ

### CLI
* Add DB to Router Creation, will Support
  * PostGres
  * SqlLite