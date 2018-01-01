package machine

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-pg/pg"
	"github.com/m49n3t0/executor/models"
	"github.com/sirupsen/logrus"
)

///////////////////////////////////////////////////////////////////////////////

// the worker executes the task process
type Worker struct {
	// ID of the local worker
	ID int64
	// Pool of worker to notify our free time
	Pool chan chan string
	// Channel of the task IDs
	Channel chan string
	// Parent dispatcher
	Dispatcher *Dispatcher
	// Quit chan
	Quit chan bool
}

// worker creation handler
func NewWorker(ID int64, dispatcher *Dispatcher) Worker {
	return Worker{
		ID:         ID,
		Pool:       dispatcher.WorkerPool,
		Channel:    make(chan string),
		Dispatcher: dispatcher,
		Quit:       make(chan bool),
	}
}

// start method starts the run loop for the worker
// listening for a quit channel in case we need to stop it
func (worker *Worker) Start() {
	go func() {

		// function logger
		log := worker.Dispatcher.Logger.
			WithField(LABEL_WORKER_ID, worker.ID)

		// infinite loop
		for {

			// register the current worker into the worker queue
			worker.Pool <- worker.Channel

			// read from the channel
			select {

			// when a task id arrived
			case ID := <-worker.Channel:

				// launch the action for this task ID
				worker.ExecuteProcess(log, ID)

			// when the quit signal arrived
			case <-worker.Quit:

				log.Info("Worker quitting")

				// we have received a signal to stop, exit this function
				return
			}
		}
	}()
}

// stop signals the worker to stop listening for work requests.
func (worker *Worker) Stop() {
	go func() {
		worker.Quit <- true
	}()
}

///////////////////////////////////////////////////////////////////////////////

// execute the action of a specific task
func (worker *Worker) ExecuteProcess(logger *logrus.Entry, ID string) {

	// retrieve a task
	task, err := worker.GetTask(ID)
	if err != nil {
		logger.WithError(err).Warn("Error while retrieve one task")
		return
	}
	// the task was already unusable, next work
	if task == nil {
		return
	}

	// function logger
	log := logger.
		WithField(LABEL_TASK_ID, task.ID).
		WithField(LABEL_STEP, task.Step)

	// logger
	log.Infof("Working on step '%s' of task '%s' ", task.Step, task.ID)

	// defined status to DOING
	task.Status = models.TaskStatus_DOING

	// update the task informations
	if err := worker.UpdateTask(task); err != nil {
		log.WithError(err).Warn("Error while update in 'DOING' status")
		return
	}

	log.Info("Task status updated in 'DOING'")

	// do the action like into a transaction
	if err := worker.DoAction(log, task); err != nil {
		// mistake strange action
		worker.ActionMistake(log, task, err)
	}

	// update the task in database
	if err := worker.UpdateTask(task); err != nil {
		log.WithError(err).Warn("Error while saving task informations")
		return
	}

	log.Info("Task correctly updated")

	return
}

// do the action for this task with the good action
func (worker *Worker) DoAction(log *logrus.Entry, task *models.Task) error {

	// get robot definition for this task version
	definition, ok := worker.Dispatcher.Definitions[task.Version]
	if !ok {
		return errors.New("Robot definition for this version doesn't exists")
	}

	// vars on step
	var step *models.Step

	// get the actual step in the sequence
	for _, s := range definition.Sequence {
		// check with the local task
		if task.Step == s.Name {
			step = &s
			break
		}
	}

	// check step found
	if step == nil {
		return errors.New("Step not found for this version")
	}

	// get the associated endpoint
	endpoint, ok := worker.Dispatcher.Endpoints[step.EndpointID]
	if !ok {
		return errors.New("Associated endpoint to this step doesn't exists")
	}

	// --------------------------------------------------------------------- //

	// do the http calls
	response, err := worker.CallHttp(log, task, endpoint)
	if err != nil {
		return worker.ActionFault(log, task, err)
	}

	// --------------------------------------------------------------------- //

	// update the local buffer from the API return
	if response.Buffer != nil {
		task.Buffer = *response.Buffer
	}

	// do the action correctly
	switch response.Action {

	// GOTO action
	case models.Action_GOTO:
		return worker.ActionGoto(log, task, definition, response)

	// GOTO_LATER action
	case models.Action_GOTO_LATER:
		return worker.ActionGotoLater(log, task, definition, response)

	// NEXT action
	case models.Action_NEXT:
		return worker.ActionNext(log, task, definition, response)

	// NEXT_LATER action
	case models.Action_NEXT_LATER:
		return worker.ActionNextLater(log, task, definition, response)

	// RETRY_NOW action
	case models.Action_RETRY_NOW:
		return worker.ActionRetryNow(log, task, definition, response)

	// RETRY action
	case models.Action_RETRY:
		return worker.ActionRetry(log, task, definition, response)

	// CANCELED action
	case models.Action_CANCELED:
		return worker.ActionCanceled(log, task, definition, response)

	// PROBLEM action
	case models.Action_PROBLEM:
		return worker.ActionProblem(log, task, definition, response)

	// ERROR action
	case models.Action_ERROR:
		return worker.ActionError(log, task, definition, response)

	}

	// action not matched
	return fmt.Errorf("Action '%s' isn't matched by executor process", response.Action)
}

///////////////////////////////////////////////////////////////////////////////

func (worker Worker) CallHttp(log *logrus.Entry, task *models.Task, endpoint *models.Endpoint) (*models.ApiResponse, error) {

	// initialize the HTTP transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// initialize the HTTP client
	httpclient := http.Client{Transport: transport}

	// parameter sends to the destination API into the request body
	var output = models.ApiParameter{
		ID:        task.ID,
		Context:   task.Context,
		Arguments: task.Arguments,
		Buffer:    task.Buffer,
	}

	// encode the body data for the call
	outputJson, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	// create the HTTP request
	request, err := http.NewRequest(endpoint.Method, endpoint.URL, bytes.NewBuffer(outputJson))
	if err != nil {
		return nil, err
	}

	// set some headers
	request.Header.Set("Content-Type", "application/json")
	//for k, v := range t.Headers {
	//	req.Header.Set(k, v)
	//}

	// timer
	start := time.Now()

	// do the HTTP call
	response, err := httpclient.Do(request)
	if err != nil {
		return nil, err
	}

	// elapsed timer
	elapsed := time.Since(start).Seconds()

	log.WithField(LABEL_HTTP_CALL_TIME, elapsed).Infof("Time of request execution : '%f' seconds", elapsed)

	// check the API return
	if response.Body == nil {
		return nil, errors.New("The API doesn't return the good structure")
	}

	// defer the closing of the body data
	defer response.Body.Close()

	// Error on the response return
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("The HTTP call failed")
	}

	// read the HTTP body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// api response
	var apiResponse = models.ApiResponse{}

	// decapsulate the body json
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, errors.New("Impossible to decapsulate the JSON body response")
	}

	return &apiResponse, nil
}

///////////////////////////////////////////////////////////////////////////////

// Function to process the GOTO_LATER action
func (worker *Worker) ActionGotoLater(log *logrus.Entry, task *models.Task, definition *models.Definition, response *models.ApiResponse) error {

	// interval settings
	//

	//default
	var interval int64 = 60

	// interval is correctly defined ?
	if response.Data.Interval != nil && *response.Data.Interval > 60 {
		interval = *response.Data.Interval
	}

	// update the task TodoDate key
	task.TodoDate = task.TodoDate.Add(time.Duration(interval) * time.Second)

	// logger
	log.Infof("TodoDate updated to '%s'", task.TodoDate) // task.TodoDate.String()

	// later == retry
	task.Retry = task.Retry - 1

	return worker.ActionGoto(log, task, definition, response)
}

// Function to process the GOTO action
func (worker *Worker) ActionGoto(log *logrus.Entry, task *models.Task, definition *models.Definition, response *models.ApiResponse) error {

	// step settings
	//

	// step is defined ?
	if response.Data.Step == nil || *response.Data.Step == "" {
		return errors.New("Missing step parameter from API response for GOTO actions")
	}

	// flag to know if found or not
	var found = false

	// asked step exists in the sequence
	for _, s := range definition.Sequence {
		// found the asked step
		if s.Name == *response.Data.Step {
			found = true
			break
		}
	}

	// not found, error
	if !found {
		return errors.New("Impossible to found the asked step from API response")
	}

	// setup the new step
	task.Step = *response.Data.Step

	// status T0D0
	task.Status = models.TaskStatus_TODO

	// logger
	log.Infof("Goto step updated to '%s'", task.Step)

	return nil
}

// Function to process the NEXT_LATER action
func (worker *Worker) ActionNextLater(log *logrus.Entry, task *models.Task, definition *models.Definition, response *models.ApiResponse) error {

	// interval settings default
	var interval int64 = 60

	// interval is correctly defined ?
	if response.Data.Interval != nil && *response.Data.Interval > 60 {
		interval = *response.Data.Interval
	}

	// update the task TodoDate key
	task.TodoDate = task.TodoDate.Add(time.Duration(interval) * time.Second)

	// logger
	log.Infof("TodoDate updated to '%s'", task.TodoDate) // task.TodoDate.String()

	// later == retry
	task.Retry = task.Retry - 1

	return worker.ActionNext(log, task, definition, response)
}

// Function to process the NEXT action
func (worker *Worker) ActionNext(log *logrus.Entry, task *models.Task, definition *models.Definition, response *models.ApiResponse) error {

	// next step settings
	//

	// flag to know if actual step found or not
	var found = false

	// store the next step name
	var nextStep string

	// retrieve the next step data
	for _, s := range definition.Sequence {
		// actual founded, this one is the classic next step
		if found {
			// next step store
			nextStep = s.Name
			break
		}
		// this actual step was here, founded
		if s.Name == task.Step {
			found = true
		}
	}

	// no next step found, error
	if !found || nextStep == "" {
		return errors.New("Impossible to found the next step")
	}

	// setup the new step
	task.Step = nextStep

	// status T0D0
	task.Status = models.TaskStatus_TODO

	// logger
	log.Infof("Next step updated to '%s'", task.Step)

	return nil
}

// Function to process the RETRY_NOW action
func (worker *Worker) ActionRetryNow(log *logrus.Entry, task *models.Task, definition *models.Definition, response *models.ApiResponse) error {

	// status T0D0
	task.Status = models.TaskStatus_TODO

	task.Retry = task.Retry - 1

	log.Info("Retry now this step")

	return nil
}

// Function to process the RETRY action
func (worker *Worker) ActionRetry(log *logrus.Entry, task *models.Task, definition *models.Definition, response *models.ApiResponse) error {

	// interval settings default
	var interval int64 = 60

	// interval is correctly defined ?
	if response.Data.Interval != nil && *response.Data.Interval > 60 {
		interval = *response.Data.Interval
	}

	// update the task TodoDate key
	task.TodoDate = task.TodoDate.Add(time.Duration(interval) * time.Second)

	// no_decrement settings
	//

	// not exists/defined no_decrement flag
	if response.Data.NoDecrement == nil || *response.Data.NoDecrement != true {
		// later == retry
		task.Retry = task.Retry - 1
	}

	// status T0D0
	task.Status = models.TaskStatus_TODO

	// logger
	log.Infof("Retry at the todoDate '%s'", task.TodoDate) // task.TodoDate.String()

	return nil
}

// Function to process the PROBLEM action
func (worker *Worker) ActionProblem(log *logrus.Entry, task *models.Task, definition *models.Definition, response *models.ApiResponse) error {

	//      comment: string               --> optional : only for ERROR/PROBLEM/CANCELED action
	//      detail: map[string]string{}   --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger

	return nil
}

// Function to process the ERROR action
func (worker *Worker) ActionError(log *logrus.Entry, task *models.Task, definition *models.Definition, response *models.ApiResponse) error {

	//      comment: string               --> optional : only for ERROR/PROBLEM/CANCELED action
	//      detail: map[string]string{}   --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger

	return nil
}

// Function to process the CANCELED action
func (worker *Worker) ActionCanceled(log *logrus.Entry, task *models.Task, definition *models.Definition, response *models.ApiResponse) error {

	//      comment: string               --> optional : only for ERROR/PROBLEM/CANCELED action
	//      detail: map[string]string{}   --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger

	return nil
}

// Function to process the MISTAKE status
func (worker *Worker) ActionMistake(log *logrus.Entry, task *models.Task, mErr error) error {

	// status MISTAKE
	task.Status = models.TaskStatus_MISTAKE

	// get the comment from error message
	task.Comment = mErr.Error()

	// always lss a retry when we do an error inside
	task.Retry = task.Retry - 1

	// logger
	log.WithError(mErr).Info("A mistake appear while the process")

	return nil
}

// Function to process the FAULT action
func (worker *Worker) ActionFault(log *logrus.Entry, task *models.Task, mErr error) error {

	// status FAULT
	task.Status = models.TaskStatus_FAULT

	// get the comment from error message
	task.Comment = mErr.Error()

	// always lss a retry when we do an error inside
	task.Retry = task.Retry - 1

	// logger
	log.WithError(mErr).Info("A disfunctional fault appear while the process")

	return nil
}

///////////////////////////////////////////////////////////////////////////////

// retrieve one task by ID in database
func (worker *Worker) GetTask(ID string) (*models.Task, error) {

	// store the data
	var task models.Task

	// fetch the object
	err := worker.Dispatcher.DB.
		Model(&task).
		Where(models.ColTask_ID+" = ?", ID).
		Where(models.ColTask_Status+" = ?", "TODO").
		Where(models.ColTask_Function+" = ?", worker.Dispatcher.Function).
		Where(models.ColTask_Retry+" > ?", 0).
		Where(models.ColTask_TodoDate + " <= NOW()").
		First()

	if err != nil {
		// return empty, no error when not found
		if err == pg.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &task, nil
}

// update the task object in database
func (worker *Worker) UpdateTask(task *models.Task) error {

	// always update the last update date key
	task.LastUpdate = time.Now()

	// update the database object
	_, err := worker.Dispatcher.DB.
		Model(task).
		Column(
			models.ColTask_Status,
			models.ColTask_Step,
			models.ColTask_Retry,
			models.ColTask_TodoDate,
			models.ColTask_Buffer,
			models.ColTask_Comment,
		).
		Update()

	if err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
