package gojenkins

import "context"

type Label struct {
	Raw     *LabelResponse
	Jenkins *Jenkins
	Base    string
}

type MODE string

const (
	NORMAL    MODE = "NORMAL"
	EXCLUSIVE      = "EXCLUSIVE"
)

type LabelNode struct {
	NodeName        string `json:"nodeName"`
	NodeDescription string `json:"nodeDescription"`
	NumExecutors    int64  `json:"numExecutors"`
	Mode            string `json:"mode"`
	Class           string `json:"_class"`
}

type LabelResponse struct {
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Nodes          []LabelNode `json:"nodes"`
	Offline        bool        `json:"offline"`
	IdleExecutors  int64       `json:"idleExecutors"`
	BusyExecutors  int64       `json:"busyExecutors"`
	TotalExecutors int64       `json:"totalExecutors"`
}

func (l *Label) GetName() string {
	return l.Raw.Name
}

func (l *Label) GetNodes() []LabelNode {
	return l.Raw.Nodes
}

func (l *Label) Poll(ctx context.Context) (int, error) {
	response, err := l.Jenkins.Requester.GetJSON(ctx, l.Base, l.Raw, nil)
	if err != nil {
		return 0, err
	}
	return response.StatusCode, nil
}
