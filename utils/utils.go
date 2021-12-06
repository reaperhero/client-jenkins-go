package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func MakeJson(data interface{}) string {
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

func ShowStatus(object string) {
	switch object {
	case "blue":
		fmt.Printf("Status: ✅ Success\n")
		break
	case "red":
		fmt.Printf("Status: ❌ Failed\n")
		break
	case "red_anime", "blue_anime", "yellow_anime", "gray_anime", "notbuild_anime":
		fmt.Printf("Status: ⏳ In Progress\n")
		break
	case "notbuilt":
		fmt.Printf("Status: 🚧 Not Build\n")
		break
	default:
		if len(object) > 0 {
			fmt.Printf("Status: %s\n", object)
		}
	}
}
