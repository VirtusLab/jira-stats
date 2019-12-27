package test

import (
	"github.com/ztrue/tracerr"
	"jira-stats/jira-stats-serverless/jira"
	"testing"
)

func TestAbs(t *testing.T) {

	_, err := jira.Fetch()
	if err != nil {
		tracerr.PrintSourceColor(err)
		t.Errorf("Found error: %s", err.Error())
	}
}
