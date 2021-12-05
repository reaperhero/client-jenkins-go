package job_test

import (
	"context"
	"fmt"
	jenkins "github.com/reaperhero/client-jenkins-go"
	"github.com/sirupsen/logrus"
	"testing"
)

type Config struct {
	Server      string
	JenkinsUser string
	Token       string
}

var (
	jc  *jenkins.Jenkins
	err error
)

func init() {
	if jc, err = jenkins.CreateJenkins(
		nil,
		"https://jenkins.mydomain.com",
		"jenkins-operator",
		"1152e8e7a88f6c7ef605844b35t5y6i",
	).Init(context.Background()); err != nil {
		panic("init client")
	}
}

func TestCreateJobInView(t *testing.T) {
	var (
		cn context.Context = context.Background()
	)
	if view, err := jc.GetView(cn, "test_list_view"); err == nil {
		ok, err := view.AddJob(cn, "jobName")
		fmt.Println(ok, err)
	}
}

func TestCreateFolder(t *testing.T) {
	var (
		cn context.Context = context.Background()
	)
	if folder,err:=jc.CreateFolder(cn,"test_folder");err==nil{
		logrus.Info(folder)
	}
}


