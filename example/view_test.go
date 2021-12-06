package example

import (
	"context"
	"fmt"
	jenkins "github.com/reaperhero/client-jenkins-go"
	"github.com/reaperhero/client-jenkins-go/utils"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

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

func TestAddJobToView(t *testing.T) {
	var (
		cn = context.Background()
	)
	if view, err := jc.GetView(cn, "test_list_view"); err == nil {
		ok, err := view.AddJob(cn, "jobName")
		fmt.Println(ok, err)
	}
}

func TestCreateFolder(t *testing.T) {
	var (
		cn = context.Background()
	)
	if folder, err := jc.CreateFolder(cn, "test_folder"); err == nil {
		logrus.Info(folder)
	}
}

func TestCreateView(t *testing.T) {
	view, err := jc.CreateView(jc.Context, "viewName", utils.DetectViewType("LIST_VIEW"))
	if err != nil {
		logrus.Println(err)
		os.Exit(1)
	}
	logrus.Info(view)
}

func TestGetPlugins(t *testing.T) {
	pls, err := jc.GetPlugins(jc.Context, 1)
	if err != nil {
		return
	}
	for _, p := range pls.Raw.Plugins {
		if len(p.LongName) > 0 && p.Active && p.Enabled {
			fmt.Printf("    %s - %s ✅\n", p.LongName, p.Version)
		}
	}
}


func TestGetView(t *testing.T) {
	views, err := jc.GetView(jc.Context, "All")
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, view := range views.Raw.Jobs {
		fmt.Printf("✅ %s\n", view.Name)
		fmt.Printf("%s\n", view.Url)
		fmt.Printf("\n")
	}
}

