package machine

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/m49n3t0/executor/models"
	"io/ioutil"
	"net/http"
	"time"
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

			case ID := <-worker.Channel:

				log.WithField(LABEL_TASK_ID, ID).Info("Worker works on task")

				// we have received a work request
				if err := worker.DoAction(ID); err != nil {
					log.WithField(LABEL_TASK_ID, ID).WithError(err).Warn("Error while execute task")
				}

			case <-worker.Quit:

				log.Info("Worker quits")

				// we have received a signal to stop
				// exit this function
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

//
//  type ApiResponse struct {
//      Action  ApiResponseAction  `json:"action,notnull"`
//      Buffer  *JsonB                  `json:"buffer"`
//      Data    *ApiResponseData   `json:"data"`
//  }
//
//  buffer: interface{},
//  action: ENUM("GOTO","NEXT","GOTO_LATER","NEXT_LATER","RETRY","RETRY_NOW","ERROR","PROBLEM","CANCELED"),
//  data: {
//      step: string                                --> optional : only for GOTO/GOTO_LATER action
//      interval: int64 in seconds ( default: 60 )  --> optional : only for GOTO_LATER/NEXT_LATER/RETRY action
//      comment: string                             --> optional : only for ERROR/PROBLEM/CANCELED action
//      detail: map[string]string{}                 --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger
//      no_decrement: bool                          --> optional : only for RETRY action
//  },
//

// do the action for this task with the good action
func (worker *Worker) DoAction(ID string) error {

	// retrieve a task
	task, err := worker.Dispatcher.GetTask(ID)
	if err != nil {
		return err
	}
    // the task was already unusable, next work
    if task == nil {
        return nil
    }

	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, ID).
        WithField(LABEL_STEP, task.Step)

	// logger
	log.Info("Working on this task")

	// --------------------------------------------------------------------- //

	// get robot for this task version
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

	// change the task data
	task.Status = models.TaskStatus_DOING

	// update the task informations
	err = worker.updateTask(task, nil)
	if err != nil {
		return err
	}

	log.Println("Task in 'DOING' status")

	// --------------------------------------------------------------------- //

	// do the http calls
	response, err := worker.CallHttp(task, endpoint)
	if err != nil {
		return worker.updateTask(task, err)
	}

	// --------------------------------------------------------------------- //

	// vars for process function response
	var mErr error

	// update the local buffer from the API return
	if response.Buffer != nil {
		task.Buffer = *response.Buffer
	}

	// do the action correctly
	switch response.Action {

	// GOTO action
	case models.Action_GOTO:
		task, mErr = worker.processActionGoto(task, definition, response)

	// GOTO_LATER action
	case models.Action_GOTO_LATER:
		task, mErr = worker.processActionGotoLater(task, definition, response)

	// NEXT action
	case models.Action_NEXT:
		task, mErr = worker.processActionNext(task, definition, response)

	// NEXT_LATER action
	case models.Action_NEXT_LATER:
		task, mErr = worker.processActionNextLater(task, definition, response)

	// RETRY_NOW action
	case models.Action_RETRY_NOW:
		task, mErr = worker.processActionRetryNow(task, definition, response)

	// RETRY action
	case models.Action_RETRY:
		task, mErr = worker.processActionRetry(task, definition, response)

	// CANCELED action
	case models.Action_CANCELED:
		task, mErr = worker.processActionCanceled(task, definition, response)

	// PROBLEM action
	case models.Action_PROBLEM:
		task, mErr = worker.processActionProblem(task, definition, response)

	// ERROR action
	case models.Action_ERROR:
		task, mErr = worker.processActionError(task, definition, response)

	// action not matched
	default:
		mErr = fmt.Errorf("Action '%s' isn't matched by executor process", response.Action)
	}

	return worker.updateTask(task, mErr)
}

func (worker *Worker) updateTask(task *models.Task, mErr error) error {

	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

	// if callback is an error, mistaken the task
	if mErr != nil {

		// status MISTAKE
		task.Status = models.TaskStatus_MISTAKE

		// get the comment from error message
		task.Comment = mErr.Error()

		// always lss a retry when we do an error inside
		task.Retry = task.Retry - 1

		log.Printf("Mistake appear : '%s'", mErr)
	}

	// update the database object
	err := worker.Dispatcher.UpdateTask(task)
	if err != nil {
		log.Printf("Error while updating the task result : %s", err)
		return err
	}

	log.Println("Task updated")

	return nil
}

///////////////////////////////////////////////////////////////////////////////

func (worker Worker) CallHttp(task *models.Task, endpoint *models.Endpoint) (*models.ApiResponse, error) {

	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

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

    log.Debugf("Time of request execution : '%f' seconds", elapsed)

	// check the API return
	if response.Body == nil {
		return nil, errors.New("The API doesn't return the good structure")
	}

	//	// XXX: need to check http status code
	//	r.StatusCode = resp.StatusCode

	// defer the closing of the body data
	defer response.Body.Close()

	// read the HTTP body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//	r.Body = string(bb)
	//
	//	bodyJSONArray := []interface{}{}
	//	if err := json.Unmarshal(bb, &bodyJSONArray); err != nil {
	//		bodyJSONMap := map[string]interface{}{}
	//		if err2 := json.Unmarshal(bb, &bodyJSONMap); err2 == nil {
	//			r.BodyJSON = bodyJSONMap
	//		}
	//	} else {
	//		r.BodyJSON = bodyJSONArray
	//	}
	fmt.Println(body)
	fmt.Println(elapsed)

	return &models.ApiResponse{}, nil
}

///////////////////////////////////////////////////////////////////////////////

// Function to process the GOTO_LATER action
func (worker *Worker) processActionGotoLater(task *models.Task, definition *models.Definition, response *models.ApiResponse) (*models.Task, error) {

	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

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
	log.Println("TodoDate updated to '%s'", task.TodoDate) // task.TodoDate.String()

	// later == retry
	task.Retry = task.Retry - 1

	return worker.processActionGoto(task, definition, response)
}

// Function to process the GOTO action
func (worker *Worker) processActionGoto(task *models.Task, definition *models.Definition, response *models.ApiResponse) (*models.Task, error) {


	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

	// step settings
	//

	// step is defined ?
	if response.Data.Step == nil || *response.Data.Step == "" {
		return task, errors.New("Missing step parameter from API response for GOTO actions")
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
		return task, errors.New("Impossible to found the asked step from API response")
	}

	// setup the new step
	task.Step = *response.Data.Step

	// status T0D0
	task.Status = models.TaskStatus_TODO

	// logger
	log.Println("Goto step updated to '%s'", task.Step)

	return task, nil
}

// Function to process the NEXT_LATER action
func (worker *Worker) processActionNextLater(task *models.Task, definition *models.Definition, response *models.ApiResponse) (*models.Task, error) {


	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

	// interval settings default
	var interval int64 = 60

	// interval is correctly defined ?
	if response.Data.Interval != nil && *response.Data.Interval > 60 {
		interval = *response.Data.Interval
	}

	// update the task TodoDate key
	task.TodoDate = task.TodoDate.Add(time.Duration(interval) * time.Second)

	// logger
	log.Println("TodoDate updated to '%s'", task.TodoDate) // task.TodoDate.String()

	// later == retry
	task.Retry = task.Retry - 1

	return worker.processActionNext(task, definition, response)
}

// Function to process the NEXT action
func (worker *Worker) processActionNext(task *models.Task, definition *models.Definition, response *models.ApiResponse) (*models.Task, error) {


	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

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
	if !found {
		return nil, errors.New("Impossible to found the next step")
	}

	// setup the new step
	task.Step = nextStep

	// status T0D0
	task.Status = models.TaskStatus_TODO

	// logger
	log.Println("Next step updated to '%s'", task.Step)

	return task, nil
}

// Function to process the RETRY_NOW action
func (worker *Worker) processActionRetryNow(task *models.Task, definition *models.Definition, response *models.ApiResponse) (*models.Task, error) {


	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

	// status T0D0
	task.Status = models.TaskStatus_TODO

	task.Retry = task.Retry - 1

	log.Println("Retry now this step")

	return task, nil
}

// Function to process the RETRY action
func (worker *Worker) processActionRetry(task *models.Task, definition *models.Definition, response *models.ApiResponse) (*models.Task, error) {


	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

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
	log.Println("Retry at the todoDate '%s'", task.TodoDate) // task.TodoDate.String()

	return task, nil
}

// Function to process the PROBLEM action
func (worker *Worker) processActionProblem(task *models.Task, definition *models.Definition, response *models.ApiResponse) (*models.Task, error) {


	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

	//      comment: string               --> optional : only for ERROR/PROBLEM/CANCELED action
	//      detail: map[string]string{}   --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger

    log.Println("toto")

	return nil, nil
}

// Function to process the ERROR action
func (worker *Worker) processActionError(task *models.Task, definition *models.Definition, response *models.ApiResponse) (*models.Task, error) {


	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)

	//      comment: string               --> optional : only for ERROR/PROBLEM/CANCELED action
	//      detail: map[string]string{}   --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger

    log.Println("toto")

	return nil, nil
}

// Function to process the CANCELED action
func (worker *Worker) processActionCanceled(task *models.Task, definition *models.Definition, response *models.ApiResponse) (*models.Task, error) {

	// function logger
	log := worker.Dispatcher.Logger.
		WithField(LABEL_WORKER_ID, worker.ID).
		WithField(LABEL_TASK_ID, task.ID).
        WithField(LABEL_STEP, task.Step)


	//      comment: string               --> optional : only for ERROR/PROBLEM/CANCELED action
	//      detail: map[string]string{}   --> optional : only for ERROR/PROBLEM/CANCELED action for push with field in the logger

    log.Println("toto")

	return nil, nil
}
