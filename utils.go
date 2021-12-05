// Copyright 2015 Vadim Kravcenko
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package gojenkins

import (
	"encoding/json"
	"fmt"
	"os"
)

func makeJson(data interface{}) string {
	str, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(json.RawMessage(str))
}


func DetectViewType(view string) string {
	viewSelected := ""
	switch view {
	case "LIST_VIEW":
		viewSelected = "hudson.model.ListView"
		break
	case "NESTED_VIEW":
		viewSelected = "hudson.plugins.nested_view.NestedView"
		break
	case "MY_VIEW":
		viewSelected = "hudson.model.MyView"
		break
	case "DASHBOARD_VIEW":
		viewSelected = "hudson.plugins.view.dashboard.Dashboard"
		break
	case "PIPELINE_VIEW":
		viewSelected = "au.com.centrumsystems.hudson.plugin.buildpipeline.BuildPipelineView"
		break
	default:
		fmt.Println("error: use only views supported: LIST_VIEW, NESTED_VIEW, MY_VIEW, DASHBOARD_VIEW, PIPELINE_VIEW")
		os.Exit(1)
	}

	return viewSelected
}