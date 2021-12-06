package gojenkins

import (
	"context"
	"fmt"
	"regexp"
)

var baseURLRegex *regexp.Regexp

func init() {
	var err error
	baseURLRegex, err = regexp.Compile("(.+)/wfapi/.*$")
	if err != nil {
		panic(err)
	}
}

type PipelineRun struct {
	Job       *Job
	Base      string
	URLs      map[string]map[string]string `json:"_links"`
	ID        string
	Name      string
	Status    string
	StartTime int64 `json:"startTimeMillis"`
	EndTime   int64 `json:"endTimeMillis"`
	Duration  int64 `json:"durationMillis"`
	Stages    []PipelineNode
}

type PipelineNode struct {
	Run            *PipelineRun
	Base           string
	URLs           map[string]map[string]string `json:"_links"`
	ID             string
	Name           string
	Status         string
	StartTime      int64 `json:"startTimeMillis"`
	Duration       int64 `json:"durationMillis"`
	StageFlowNodes []PipelineNode
	ParentNodes    []int64
}

type PipelineInputAction struct {
	ID         string
	Message    string
	ProceedURL string
	AbortURL   string
}

type PipelineArtifact struct {
	ID   string
	Name string
	Path string
	URL  string
	size int
}

type PipelineNodeLog struct {
	NodeID     string
	NodeStatus string
	Length     int64
	HasMore    bool
	Text       string
	ConsoleURL string
}

func (run *PipelineRun) update() {
	href := run.URLs["self"]["href"]
	if matches := baseURLRegex.FindStringSubmatch(href); len(matches) > 1 {
		run.Base = matches[1]
	}
	for i := range run.Stages {
		run.Stages[i].Run = run
		href := run.Stages[i].URLs["self"]["href"]
		if matches := baseURLRegex.FindStringSubmatch(href); len(matches) > 1 {
			run.Stages[i].Base = matches[1]
		}
	}
}

func (j *Job) GetPipelineRuns(ctx context.Context) (pr []PipelineRun, err error) {
	_, err = job.Jenkins.Requester.GetJSON(ctx, job.Base+"/wfapi/runs", &pr, nil)
	if err != nil {
		return nil, err
	}
	for i := range pr {
		pr[i].update()
		pr[i].Job = job
	}

	return pr, nil
}

func (j *Job) GetPipelineRun(ctx context.Context, id string) (pr *PipelineRun, err error) {
	pr = new(PipelineRun)
	href := job.Base + "/" + id + "/wfapi/describe"
	_, err = job.Jenkins.Requester.GetJSON(ctx, href, pr, nil)
	if err != nil {
		return nil, err
	}
	pr.update()
	pr.Job = job

	return pr, nil
}

func (run *PipelineRun) GetPendingInputActions(ctx context.Context) (PIAs []PipelineInputAction, err error) {
	PIAs = make([]PipelineInputAction, 0, 1)
	href := run.Base + "/wfapi/pendingInputActions"
	_, err = run.Job.Jenkins.Requester.GetJSON(ctx, href, &PIAs, nil)
	if err != nil {
		return nil, err
	}

	return PIAs, nil
}

func (run *PipelineRun) GetArtifacts(ctx context.Context) (artifacts []PipelineArtifact, err error) {
	artifacts = make([]PipelineArtifact, 0, 0)
	href := run.Base + "/wfapi/artifacts"
	_, err = run.Job.Jenkins.Requester.GetJSON(ctx, href, artifacts, nil)
	if err != nil {
		return nil, err
	}

	return artifacts, nil
}

func (run *PipelineRun) GetNode(ctx context.Context, id string) (node *PipelineNode, err error) {
	node = new(PipelineNode)
	href := run.Base + "/execution/node/" + id + "/wfapi/describe"
	_, err = run.Job.Jenkins.Requester.GetJSON(ctx, href, node, nil)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (node *PipelineNode) GetLog(ctx context.Context) (log *PipelineNodeLog, err error) {
	log = new(PipelineNodeLog)
	href := node.Base + "/wfapi/log"
	fmt.Println(href)
	_, err = node.Run.Job.Jenkins.Requester.GetJSON(ctx, href, log, nil)
	if err != nil {
		return nil, err
	}

	return log, nil
}
