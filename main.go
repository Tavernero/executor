package executor

import (
	"fmt"
    "strconv"
)

type Task struct {
    ID int

    Function string

    Name string

    Step string
    Status string
    Retry int

    Arguments string
    Buffer string
}

type Step struct {
    Name string
    Url string
}

type Configuration struct {
    Name string
    Steps []Step
}

type Robot struct {
    Conf Configuration
    Debug bool
}

// Create a new Executor step-machine robot
func New(function string, debug bool) (bot *Robot, err error) {

    // XXX : default sequence for testing period
    // XXX   need to retrieve the configuration from a database
    configuration := Configuration{
        Name: function,
        Steps: []Step{
            Step{
                Name:"starting",
                Url:"api.com/starting",
            },
            Step{
                Name:"onServer",
                Url:"api.com/on_server",
            },
            Step{
                Name:"onInterne",
                Url:"api.com/on_interne",
            },
            Step{
                Name:"ending",
                Url:"api.com/ending",
            },
        },
    }
    // XXX

    bot = &Robot{
        Conf: configuration,
        Debug: debug,
    }

    if bot.Debug {
        fmt.Println("==============================================================")
        fmt.Println( bot )
        fmt.Println("==============================================================")
    }

    return
}

// Launch the executor run
func (bot *Robot) Run() (err error) {

    // XXX : need to fetch the task from database
    tasks := []Task{
        Task{
            ID: 3,
            Function: "database/create",
            Name: "toto",
            Step: "starting",
            Status: "todo",
            Retry: 8,
            Arguments: "{'arguments':{}}",
            Buffer: "{'buffer':'{}}",
        },
        Task{
            ID: 7,
            Function: "database/create",
            Name: "zozo",
            Step: "starting",
            Status: "todo",
            Retry: 8,
            Arguments: "{'arguments':{}}",
            Buffer: "{'buffer':'{}}",
        },
        Task{
            ID: 9,
            Function: "database/create",
            Name: "popo",
            Step: "starting",
            Status: "todo",
            Retry: 8,
            Arguments: "{'arguments':{}}",
            Buffer: "{'buffer':'{}}",
        },
    }
    // XXX

    if bot.Debug {
        fmt.Println("==============================================================")
        fmt.Println( tasks )
        fmt.Println("==============================================================")
    }


    for task_num := range tasks {
        task := tasks[task_num]
        prefix := "[Task:" + strconv.Itoa(task.ID) + "] " + task.Name + " -- " + " "

        fmt.Print( prefix )
        fmt.Println("Begin")

        for num := range bot.Conf.Steps {
            step := bot.Conf.Steps[num]

            fmt.Print( prefix )
            fmt.Print("Do step '")
            fmt.Print( step.Name )
            fmt.Print("' calling '")
            fmt.Print( step.Url )
            fmt.Println("'")
        }

        fmt.Print( prefix )
        fmt.Println("End")

        fmt.Println(" --- next ---")
    }

	return err
}
