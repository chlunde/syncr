package main

import (
	"fmt"
	"io"

	"github.com/chlunde/syncr"
)

const (
	ClearAndMoveHome = "\033[2J\033[2H"
	CyanBold         = "\033[36;1m"
	RedBold          = "\033[31;1m"
	GreenBold        = "\033[32;1m"
	ResetColor       = "\033[0m"
)

func Display(buf io.Writer, syncrs []*syncr.Syncr) bool {
	fmt.Fprint(buf, ClearAndMoveHome)

	var anyAlive bool
	for _, s := range syncrs {
		s.Lock.Lock()
		fmt.Fprint(buf, CyanBold, s.Description)

		var statusColor string
		var sep string
		switch {
		case s.Error != nil:
			statusColor = RedBold
			sep = "\n"
		case s.Dead:
			statusColor = GreenBold
			sep = " "
		default:
			statusColor = ResetColor
			sep = "\n"
		}

		fmt.Fprint(buf, statusColor, sep)

		if s.Error != nil {
			fmt.Fprintf(buf, "Failed: %v\n", s.Error)
		}

		fmt.Fprintln(buf, s.Status[0])

		// Skip last line if it's blank we're not getting any more lines
		// prevents the status from jumping up and down
		if !s.Dead || len(s.Status[1]) > 0 {
			fmt.Fprintln(buf, s.Status[1])
		}

		fmt.Fprint(buf, ResetColor)

		if !s.Dead {
			anyAlive = true
		}
		s.Lock.Unlock()
	}

	return anyAlive
}
