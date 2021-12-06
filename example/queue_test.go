package example

import (
	"fmt"
	"github.com/reaperhero/client-jenkins-go/utils"
	"testing"
)

func TestShowBuildQueue(t *testing.T) {
	queue, _ := jc.GetQueue(jc.Context)
	totalTasks := 0
	for i, item := range queue.Raw.Items {
		fmt.Printf("Name: %s\n", item.Task.Name)
		fmt.Printf("ID: %d\n", item.ID)
		utils.ShowStatus(item.Task.Color)
		fmt.Printf("Pending: %v\n", item.Pending)
		fmt.Printf("Stuck: %v\n", item.Stuck)

		fmt.Printf("Why: %s\n", item.Why)
		fmt.Printf("URL: %s\n", item.Task.URL)
		fmt.Printf("\n")
		totalTasks = i + 1
	}
	fmt.Printf("Number of tasks in the build queue: %d\n", totalTasks)
}
