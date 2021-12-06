package example

import (
	"fmt"
	"github.com/reaperhero/client-jenkins-go/utils"
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

func TestEnableJob(t *testing.T) {
	job, err := jc.GetJob(jc.Context, "job_id")
	if err != nil {
		os.Exit(1)
	}
	_, _ = job.Enable(jc.Context)
}

func TestGetLastUnstableBuild(t *testing.T) {
	job, err := jc.GetJob(jc.Context, "jobName")
	if err != nil {
		logrus.Info("❌ unable to find the specific job")
	}
	build, err := job.GetLastBuild(jc.Context)
	if err != nil {
		logrus.Info("❌ unable to find the last unstable build job")
	}

	if len(build.Job.Raw.LastBuild.URL) > 0 {
		fmt.Printf("Last unstable build Number: %d\n", build.Job.Raw.LastBuild.Number)
		fmt.Printf("Last unstable build URL: %s\n", build.Job.Raw.LastBuild.URL)
		fmt.Printf("Parameters: %s\n", build.GetParameters())
	} else {
		fmt.Printf("No last unstable build available for job: %s", "jobName")
	}
}

func TestShowAllJobs(t *testing.T) {
	jobs, err := jc.GetAllJobs(jc.Context)
	if err != nil {
		logrus.Error(err)
	}
	for _, job := range jobs {
		fmt.Printf("✅ %s\n", job.Raw.Name)
		utils.ShowStatus(job.Raw.Color)
		fmt.Printf("%s\n", job.Raw.Description)
		fmt.Printf("%s\n", job.Raw.URL)
		fmt.Printf("\n")
	}
}
