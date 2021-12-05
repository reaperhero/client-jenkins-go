package example

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestCreateNode(t *testing.T) {
	node, err := jc.CreateNode(
		jc.Context,
		"NODE_NAME",
		10,
		"DESCRIPTION",
		"REMOTEFS",
		"LABEL",
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logrus.Info(node)
}

func TestDeleteNode(t *testing.T) {
	_, _ = jc.DeleteNode(jc.Context, "nodeName")
}
