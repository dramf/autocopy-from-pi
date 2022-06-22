package main

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		title     string
		version   string
		buildtime string
		result    string
	}{
		{
			title:     "good way",
			version:   "1.1.1",
			buildtime: "1655838150",
			result:    "v1.1.1 build: 2022.06.22 01:02:30",
		},
		{
			title:   "without buildtime",
			version: "0.1.1",
			result:  "v0.1.1 build: 2022.06.16 14:12:43",
		},
		{
			title:     "without version",
			buildtime: "1655838150",
			result:    "v build: 2022.06.22 01:02:30",
		},
	}

	defaultVersion := version
	defaultBuildtime := buildtime
	defer func() {
		version = defaultVersion
		buildtime = defaultBuildtime
	}()

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			version = test.version
			buildtime = test.buildtime
			got := getVersion()
			if got != test.result {
				t.Errorf("getVersion(%q,%q) == %q, wanted %q", test.version, test.buildtime, got, test.result)
			}
		})
	}
}
