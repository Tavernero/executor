package executor

import (
	"fmt"
)

// Bot contains informations for current robots skeleton
type Bot struct {
    function string
}

func New(function string) (*Bot, error) {

    var bot Bot
    var err error

    bot = Bot{
        function: function,
    }

    fmt.Println( bot )

    return &bot, err
}
