package gojenkins

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/reaperhero/client-jenkins-go/utils"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type Job struct {
	Raw     *JobResponse
	Jenkins *Jenkins
	Base    string
}

type JobBuild struct {
	Number int64
	URL    string
}

type InnerJob struct {
	Class string `json:"_class"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Color string `json:"color"`
}

type ParameterDefinition struct {
	DefaultParameterValue struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"defaultParameterValue"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type JobResponse struct {
	Class              string `json:"_class"`
	Actions            []generalObj
	Buildable          bool `json:"buildable"`
	Builds             []JobBuild
	Color              string      `json:"color"`
	ConcurrentBuild    bool        `json:"concurrentBuild"`
	Description        string      `json:"description"`
	DisplayName        string      `json:"displayName"`
	DisplayNameOrNull  interface{} `json:"displayNameOrNull"`
	DownstreamProjects []InnerJob  `json:"downstreamProjects"`
	FirstBuild         JobBuild
	FullName           string `json:"fullName"`
	FullDisplayName    string `json:"fullDisplayName"`
	HealthReport       []struct {
		Description   string `json:"description"`
		IconClassName string `json:"iconClassName"`
		IconUrl       string `json:"iconUrl"`
		Score         int64  `json:"score"`
	} `json:"healthReport"`
	InQueue               bool     `json:"inQueue"`
	KeepDependencies      bool     `json:"keepDependencies"`
	LastBuild             JobBuild `json:"lastBuild"`
	LastCompletedBuild    JobBuild `json:"lastCompletedBuild"`
	LastFailedBuild       JobBuild `json:"lastFailedBuild"`
	LastStableBuild       JobBuild `json:"lastStableBuild"`
	LastSuccessfulBuild   JobBuild `json:"lastSuccessfulBuild"`
	LastUnstableBuild     JobBuild `json:"lastUnstableBuild"`
	LastUnsuccessfulBuild JobBuild `json:"lastUnsuccessfulBuild"`
	Name                  string   `json:"name"`
	NextBuildNumber       int64    `json:"nextBuildNumber"`
	Property              []struct {
		ParameterDefinitions []ParameterDefinition `json:"parameterDefinitions"`
	} `json:"property"`
	QueueItem        interface{} `json:"queueItem"`
	Scm              struct{}    `json:"scm"`
	UpstreamProjects []InnerJob  `json:"upstreamProjects"`
	URL              string      `json:"url"`
	Jobs             []InnerJob  `json:"jobs"`
	PrimaryView      *ViewData   `json:"primaryView"`
	Views            []ViewData  `json:"views"`
}

func (job *Job) parentBase() string {
	return job.Base[:strings.LastIndex(job.Base, "/job/")]
}

type History struct {
	BuildDisplayName string
	BuildNumber      int
	BuildStatus      string
	BuildTimestamp   int64
}

func (job *Job) GetName() string {
	return job.Raw.Name
}

func (job *Job) GetDescription() string {
	return job.Raw.Description
}

func (job *Job) GetDetails() *JobResponse {
	return job.Raw
}

func (job *Job) GetBuild(ctx context.Context, id int64) (*Build, error) {

	// Support customized server URL,
	// i.e. Server : https://<domain>/jenkins/job/JOB1
	// "https://<domain>/jenkins/" is the server URL,
	// we are expecting jobURL = "job/JOB1"
	jobURL := strings.Replace(job.Raw.URL, job.Jenkins.Server, "", -1)
	build := Build{Jenkins: job.Jenkins, Job: job, Raw: new(BuildResponse), Depth: 1, Base: jobURL + "/" + strconv.FormatInt(id, 10)}
	status, err := build.Poll(ctx)
	if err != nil {
		return nil, err
	}
	if status == 200 {
		return &build, nil
	}
	return nil, errors.New(strconv.Itoa(status))
}

func (job *Job) getBuildByType(ctx context.Context, buildType string) (*Build, error) {
	allowed := map[string]JobBuild{
		"lastStableBuild":     job.Raw.LastStableBuild,
		"lastSuccessfulBuild": job.Raw.LastSuccessfulBuild,
		"lastBuild":           job.Raw.LastBuild,
		"lastCompletedBuild":  job.Raw.LastCompletedBuild,
		"firstBuild":          job.Raw.FirstBuild,
		"lastFailedBuild":     job.Raw.LastFailedBuild,
	}
	number := ""
	if val, ok := allowed[buildType]; ok {
		number = strconv.FormatInt(val.Number, 10)
	} else {
		panic("No Such Build")
	}
	build := Build{
		Jenkins: job.Jenkins,
		Depth:   1,
		Job:     job,
		Raw:     new(BuildResponse),
		Base:    job.Base + "/" + number}
	status, err := build.Poll(ctx)
	if err != nil {
		return nil, err
	}
	if status == 200 {
		return &build, nil
	}
	return nil, errors.New(strconv.Itoa(status))
}

func (job *Job) GetLastSuccessfulBuild(ctx context.Context) (*Build, error) {
	return job.getBuildByType(ctx, "lastSuccessfulBuild")
}

func (job *Job) GetFirstBuild(ctx context.Context) (*Build, error) {
	return job.getBuildByType(ctx, "firstBuild")
}

func (job *Job) GetLastBuild(ctx context.Context) (*Build, error) {
	return job.getBuildByType(ctx, "lastBuild")
}

func (job *Job) GetLastStableBuild(ctx context.Context) (*Build, error) {
	return job.getBuildByType(ctx, "lastStableBuild")
}

func (job *Job) GetLastFailedBuild(ctx context.Context) (*Build, error) {
	return job.getBuildByType(ctx, "lastFailedBuild")
}

func (job *Job) GetLastCompletedBuild(ctx context.Context) (*Build, error) {
	return job.getBuildByType(ctx, "lastCompletedBuild")
}

func (job *Job) GetBuildsFields(ctx context.Context, fields []string, custom interface{}) error {
	if fields == nil || len(fields) == 0 {
		return fmt.Errorf("one or more field value needs to be specified")
	}
	// limit overhead using builds instead of allBuilds, which returns the last 100 build
	_, err := job.Jenkins.Requester.GetJSON(ctx, job.Base, &custom, map[string]string{"tree": "builds[" + strings.Join(fields, ",") + "]"})
	if err != nil {
		return err
	}
	return nil
}

// Returns All Builds with Number and URL
func (job *Job) GetAllBuildIds(ctx context.Context) ([]JobBuild, error) {
	var buildsResp struct {
		Builds []JobBuild `json:"allBuilds"`
	}
	_, err := job.Jenkins.Requester.GetJSON(ctx, job.Base, &buildsResp, map[string]string{"tree": "allBuilds[number,url]"})
	if err != nil {
		return nil, err
	}
	return buildsResp.Builds, nil
}

func (job *Job) GetUpstreamJobsMetadata() []InnerJob {
	return job.Raw.UpstreamProjects
}

func (job *Job) GetDownstreamJobsMetadata() []InnerJob {
	return job.Raw.DownstreamProjects
}

func (job *Job) GetInnerJobsMetadata() []InnerJob {
	return job.Raw.Jobs
}

func (job *Job) GetUpstreamJobs(ctx context.Context) ([]*Job, error) {
	jobs := make([]*Job, len(job.Raw.UpstreamProjects))
	for i, job := range job.Raw.UpstreamProjects {
		ji, err := job.Jenkins.GetJob(ctx, job.Name)
		if err != nil {
			return nil, err
		}
		jobs[i] = ji
	}
	return jobs, nil
}

func (job *Job) GetDownstreamJobs(ctx context.Context) ([]*Job, error) {
	jobs := make([]*Job, len(job.Raw.DownstreamProjects))
	for i, job := range job.Raw.DownstreamProjects {
		ji, err := job.Jenkins.GetJob(ctx, job.Name)
		if err != nil {
			return nil, err
		}
		jobs[i] = ji
	}
	return jobs, nil
}

func (job *Job) GetInnerJob(ctx context.Context, id string) (*Job, error) {
	job := Job{Jenkins: job.Jenkins, Raw: new(JobResponse), Base: job.Base + "/job/" + id}
	status, err := job.Poll(ctx)
	if err != nil {
		return nil, err
	}
	if status == 200 {
		return &job, nil
	}
	return nil, errors.New(strconv.Itoa(status))
}

func (job *Job) GetInnerJobs(ctx context.Context) ([]*Job, error) {
	jobs := make([]*Job, len(job.Raw.Jobs))
	for i, job := range job.Raw.Jobs {
		ji, err := job.GetInnerJob(ctx, job.Name)
		if err != nil {
			return nil, err
		}
		jobs[i] = ji
	}
	return jobs, nil
}

func (job *Job) Enable(ctx context.Context) (bool, error) {
	resp, err := job.Jenkins.Requester.Post(ctx, job.Base+"/enable", nil, nil, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New(strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}

func (job *Job) Disable(ctx context.Context) (bool, error) {
	resp, err := job.Jenkins.Requester.Post(ctx, job.Base+"/disable", nil, nil, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New(strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}

func (job *Job) Delete(ctx context.Context) (bool, error) {
	resp, err := job.Jenkins.Requester.Post(ctx, job.Base+"/doDelete", nil, nil, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New(strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}

func (job *Job) Rename(ctx context.Context, name string) (bool, error) {
	data := url.Values{}
	data.Set("newName", name)
	_, err := job.Jenkins.Requester.Post(ctx, job.Base+"/doRename", bytes.NewBufferString(data.Encode()), nil, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (job *Job) Create(ctx context.Context, config string, qr ...interface{}) (*Job, error) {
	var querystring map[string]string
	if len(qr) > 0 {
		querystring = qr[0].(map[string]string)
	}
	resp, err := job.Jenkins.Requester.PostXML(ctx, job.parentBase()+"/createItem", config, job.Raw, querystring)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		job.Poll(ctx)
		return job, nil
	}
	return nil, errors.New(strconv.Itoa(resp.StatusCode))
}

func (job *Job) Copy(ctx context.Context, destinationName string) (*Job, error) {
	qr := map[string]string{"name": destinationName, "from": job.GetName(), "mode": "copy"}
	resp, err := job.Jenkins.Requester.Post(ctx, job.parentBase()+"/createItem", nil, nil, qr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		newJob := &Job{Jenkins: job.Jenkins, Raw: new(JobResponse), Base: "/job/" + destinationName}
		_, err := newJob.Poll(ctx)
		if err != nil {
			return nil, err
		}
		return newJob, nil
	}
	return nil, errors.New(strconv.Itoa(resp.StatusCode))
}

func (job *Job) UpdateConfig(ctx context.Context, config string) error {

	var querystring map[string]string

	resp, err := job.Jenkins.Requester.PostXML(ctx, job.Base+"/config.xml", config, nil, querystring)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		job.Poll(ctx)
		return nil
	}
	return errors.New(strconv.Itoa(resp.StatusCode))

}

func (job *Job) GetConfig(ctx context.Context) (string, error) {
	var data string
	_, err := job.Jenkins.Requester.GetXML(ctx, job.Base+"/config.xml", &data, nil)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (job *Job) GetParameters(ctx context.Context) ([]ParameterDefinition, error) {
	_, err := job.Poll(ctx)
	if err != nil {
		return nil, err
	}
	var parameters []ParameterDefinition
	for _, property := range job.Raw.Property {
		parameters = append(parameters, property.ParameterDefinitions...)
	}
	return parameters, nil
}

func (job *Job) IsQueued(ctx context.Context) (bool, error) {
	if _, err := job.Poll(ctx); err != nil {
		return false, err
	}
	return job.Raw.InQueue, nil
}

func (job *Job) IsRunning(ctx context.Context) (bool, error) {
	if _, err := job.Poll(ctx); err != nil {
		return false, err
	}
	lastBuild, err := job.GetLastBuild(ctx)
	if err != nil {
		return false, err
	}
	return lastBuild.IsRunning(ctx), nil
}

func (job *Job) IsEnabled(ctx context.Context) (bool, error) {
	if _, err := job.Poll(ctx); err != nil {
		return false, err
	}
	return job.Raw.Color != "disabled", nil
}

func (job *Job) HasQueuedBuild() {
	panic("Not Implemented yet")
}

func (job *Job) InvokeSimple(ctx context.Context, params map[string]string) (int64, error) {
	isQueued, err := job.IsQueued(ctx)
	if err != nil {
		return 0, err
	}
	if isQueued {
		Error.Printf("%s is already running", job.GetName())
		return 0, nil
	}

	endpoint := "/build"
	parameters, err := job.GetParameters(ctx)
	if err != nil {
		return 0, err
	}
	if len(parameters) > 0 {
		endpoint = "/buildWithParameters"
	}
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}
	resp, err := job.Jenkins.Requester.Post(ctx, job.Base+endpoint, bytes.NewBufferString(data.Encode()), nil, nil)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return 0, fmt.Errorf("Could not invoke job %q: %s", job.GetName(), resp.Status)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return 0, errors.New("Don't have key \"Location\" in response of header")
	}

	u, err := url.Parse(location)
	if err != nil {
		return 0, err
	}

	number, err := strconv.ParseInt(path.Base(u.Path), 10, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func (job *Job) Invoke(ctx context.Context, files []string, skipIfRunning bool, params map[string]string, cause string, securityToken string) (bool, error) {
	isQueued, err := job.IsQueued(ctx)
	if err != nil {
		return false, err
	}
	if isQueued {
		Error.Printf("%s is already running", job.GetName())
		return false, nil
	}
	isRunning, err := job.IsRunning(ctx)
	if err != nil {
		return false, err
	}
	if isRunning && skipIfRunning {
		return false, fmt.Errorf("Will not request new build because %s is already running", job.GetName())
	}

	base := "/build"

	// If parameters are specified - url is /builWithParameters
	if params != nil {
		base = "/buildWithParameters"
	} else {
		params = make(map[string]string)
	}

	// If files are specified - url is /build
	if files != nil {
		base = "/build"
	}
	reqParams := map[string]string{}
	buildParams := map[string]string{}
	if securityToken != "" {
		reqParams["token"] = securityToken
	}

	buildParams["json"] = string(utils.MakeJson(params))
	b, _ := json.Marshal(buildParams)
	resp, err := job.Jenkins.Requester.PostFiles(ctx, job.Base+base, bytes.NewBuffer(b), nil, reqParams, files)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return true, nil
	}
	return false, errors.New(strconv.Itoa(resp.StatusCode))
}

func (job *Job) Poll(ctx context.Context) (int, error) {
	response, err := job.Jenkins.Requester.GetJSON(ctx, job.Base, job.Raw, nil)
	if err != nil {
		return 0, err
	}
	return response.StatusCode, nil
}

func (job *Job) History(ctx context.Context) ([]*History, error) {
	var s string
	_, err := job.Jenkins.Requester.Get(ctx, job.Base+"/buildHistory/ajax", &s, nil)
	if err != nil {
		return nil, err
	}

	return parseBuildHistory(strings.NewReader(s)), nil
}

func (run *PipelineRun) ProceedInput(ctx context.Context) (bool, error) {
	actions, _ := run.GetPendingInputActions(ctx)
	data := url.Values{}
	data.Set("inputId", actions[0].ID)
	params := make(map[string]string)
	data.Set("json", utils.MakeJson(params))

	href := run.Base + "/wfapi/inputSubmit"

	resp, err := run.Job.Jenkins.Requester.Post(ctx, href, bytes.NewBufferString(data.Encode()), nil, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New(strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}

func (run *PipelineRun) AbortInput(ctx context.Context) (bool, error) {
	actions, _ := run.GetPendingInputActions(ctx)
	data := url.Values{}
	params := make(map[string]string)
	data.Set("json", utils.MakeJson(params))

	href := run.Base + "/input/" + actions[0].ID + "/abort"

	resp, err := run.Job.Jenkins.Requester.Post(ctx, href, bytes.NewBufferString(data.Encode()), nil, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New(strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}
