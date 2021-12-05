package example

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestCreateJobInFolder(t *testing.T) {
	job, err := jc.CreateJobInFolder(jc.Context, "<config></config>", "newJobName", "myFolder", "parentFolder") //JOB_DATA_XML JOB_NAME FOLDER_NAME
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	fmt.Println(job)
}

func TestCreateJob(t *testing.T) {
	job, err := jc.CreateJob(jc.Context, "job.xml", "jobName")
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(job)

}

func TestDeleteJob(t *testing.T) {
	_, _ = jc.DeleteJob(jc.Context, "jobName")
}
