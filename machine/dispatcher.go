package machine

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
    "time"

	"github.com/go-pg/pg"
	"github.com/m49n3t0/executor/models"
	"github.com/sirupsen/logrus"
)

///////////////////////////////////////////////////////////////////////////////

// to define the dispatcher parameters
type DispatcherParams struct {
	// robot configuration
	Function string
	// dispatcher configuration
	MaxWorker int64
	MaxQueue  int64
	// machine interface object
	Machine Machine
}

// dispatcher object
type Dispatcher struct {
	// function on what the robot work
	Function string
	// store datas from database
	Definitions map[int64]*models.Definition
	Endpoints   map[string]*models.Endpoint
	// manage the distribution of workflow
	WorkerPool chan chan string
	Queue      chan string
	MaxWorker  int64
	MaxQueue   int64
	// manage the quit process
	Signals chan os.Signal
	Quit    chan bool
	// database handler
	DB *pg.DB
	// logger handler
	Logger *logrus.Entry
}

// dispatcher creation handler
func NewDispatcher(params *DispatcherParams) (*Dispatcher, error) {

	// machine interface object
	machine := params.Machine

	// get the logger from interface
	logger, err := machine.GetLogger()
	if err != nil {
		return nil, err
	}
	if logger == nil {
		return nil, errors.New("the logger initialization return an empty object")
	}

	// encapsulate function field
	logentry := logger.WithField(LABEL_FUNCTION, params.Function)

	// get the database handler from interface
	database, err := machine.GetDatabase(logentry)
	if err != nil {
		return nil, err
	}

	// create the object
	dispatcher := &Dispatcher{
		Function:    params.Function,
		Definitions: make(map[int64]*models.Definition),
		Endpoints:   make(map[string]*models.Endpoint),
		WorkerPool:  make(chan chan string, params.MaxWorker),
		Queue:       make(chan string, params.MaxQueue),
		MaxWorker:   params.MaxWorker,
		MaxQueue:    params.MaxQueue,
		Signals:     make(chan os.Signal, 2),
		Quit:        make(chan bool),
		DB:          database,
		Logger:      logentry,
	}

	return dispatcher, nil
}

// stop signals programmatically
func (dispatcher *Dispatcher) Stop() {
	go func() {
		dispatcher.Quit <- true
	}()
}

// stop signals from system
func (dispatcher *Dispatcher) Signal() {
	go func() {

		// link system signal to the dispatcher signal
		signal.Notify(dispatcher.Signals, os.Interrupt, syscall.SIGTERM)

		// when receive syscall signal
		<-dispatcher.Signals

		// do a stopper
		dispatcher.Stop()
	}()
}

///////////////////////////////////////////////////////////////////////////////

// retrieve robot configuration for this function from database
func (dispatcher *Dispatcher) GetRobotConfiguration() error {

	// function logger
	log := dispatcher.Logger

	log.Info("Get the robot configuration")

	// object to fetch
	var robots []*models.Robot

	// get the robot data
	err := dispatcher.DB.
		Model(&robots).
		Where(models.ColRobot_Function+" = ?", dispatcher.Function).
		Where(models.ColRobot_Status+" = ?", "ACTIVE").
		Select()

	if err != nil {
		log.WithError(err).Warn("Select robot configuration error")
		return err
	}

	// no elements, error
	if len(robots) == 0 {
		return errors.New("No robots definition active found for this function")
	}

	// store the endpoint to fetch informations
	var stepIDs []string

    // store into dispatcher definitions data
	for _, robot := range robots {

		// remap by version the definitions
		dispatcher.Definitions[robot.Version] = &robot.Definition

		// save the step IDs to fetch after
		for _, step := range robot.Definition.Sequence {
			stepIDs = append(stepIDs, step.EndpointID)
		}
	}

	// no step elements, error
	if len(stepIDs) == 0 {
		return errors.New("No steps defined into action robots")
	}

	// object to fetch
	var endpoints []*models.Endpoint

	// get the endpoint data
	err = dispatcher.DB.
		Model(&endpoints).
		Where(models.ColEndpoint_ID+" IN ( ? )", pg.In(stepIDs)).
		Select()

	if err != nil {
		log.WithError(err).Warn("Select robot endpoint steps error")
		return err
	}

	// store locally each endpoints
	for _, endpoint := range endpoints {
		dispatcher.Endpoints[endpoint.ID] = endpoint
	}

	log.Info("Robot configuration loaded")

	return nil
}

// retrieve the available task ID list
func (dispatcher *Dispatcher) GetTasks() error {

	// function logger
	log := dispatcher.Logger

	log.Info("Reading task ID")

	// where store the ID list
	var IDs []string

	// working on this model
	var task models.Task

	// fetch the available ID list
	err := dispatcher.DB.
		Model(&task).
		Column(models.ColTask_ID).
		Where(models.ColTask_Status+" = ?", models.TaskStatus_TODO).
		Where(models.ColTask_Function+" = ?", dispatcher.Function).
		Where(models.ColTask_Retry+" > ?", 0).
		Where(models.ColTask_TodoDate + " <= NOW()").
		OrderExpr(models.ColTask_TodoDate + " ASC").
		Select(&IDs)

	if err != nil {
		log.WithError(err).Warn("Select task IDs error")
		return err
	}

	// push in the queue the ID informations
	for _, ID := range IDs {
		dispatcher.Queue <- ID
	}

	return nil
}

// retrieve one task by ID
func (dispatcher *Dispatcher) GetTask(ID string) (*models.Task, error) {

	// store the data
	var task models.Task

	// fetch the object
	err := dispatcher.DB.
		Model(&task).
		Where(models.ColTask_ID+" = ?", ID).
		Where(models.ColTask_Status+" = ?", "TODO").
		Where(models.ColTask_Function+" = ?", dispatcher.Function).
		Where(models.ColTask_Retry+" > ?", 0).
		Where(models.ColTask_TodoDate + " <= NOW()").
		First()

	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("Task not found")
		}

		return nil, err
	}

	return &task, nil
}

// update the task object in database
func (dispatcher *Dispatcher) UpdateTask(task *models.Task) error {

	// function logger
	log := dispatcher.Logger

	// always update the last update date key
	task.LastUpdate = time.Now()

	// update the database object
    _, err := dispatcher.DB.
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
		log.Printf("Error while updating the task result : %s", err)
		return err
	}

	log.Println("Task updated")

	return nil
}


