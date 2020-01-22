# jira-stats
[![Travis CI](https://img.shields.io/travis/VirtusLab/jira-stats.svg)](https://travis-ci.org/VirtusLab/render)

Simple tool to fetching Jira tickets and genereting CSV for reporting purposes.

#### Building codebase locally:
* Make sure you have [Golang](https://golang.org/doc/install) installed
* Make sure you have [Serverless Framework](https://serverless.com/framework/docs/getting-started/) tools installed  
* Create dir `$GO_PATH/github.com/VirtusLab` and clone this repo there
* Go into repo dir and type `make build`


#### Running it locally
* Make sure to export two env variables:

        JIRA_USER
        JIRA_PASSWORD 

* run `local/main.go` - if running from IDE make sure that env variables are visible there

#### To deploy
* Make sure you have JIRA env vars exported (look above)
* Run: `sls deploy`