package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/leonidasdeim/backupper/internal"
)

type FilterState interface {
	SetFilter(filter internal.Filter) *internal.State
	GetFilter() *internal.Filter
	Save() error
}

type viewer struct {
	log    internal.Log
	filter internal.Filter
}

var utils = internal.Utils

// Runs log filter prompt setup and log viewer runner.
// This function is blocking, should be run from goroutine.
func RunLogViewer(log internal.Log, state FilterState) {
	if state == nil {
		fmt.Println("Can't run log viewer: state object is not provided")
		return
	}

	if proceed := utils.Prompt("Do you want to view log? y/n"); proceed != "y" {
		fmt.Println("Log viewer disabled")
		return
	}

	viewer{
		log:    log,
		filter: configureFilter(state),
	}.runner()
}

func (lf viewer) runner() {
	fmt.Println("Log viewer started...")

	utils.FileScanner(
		lf.log.GetFile(),
		func(line string) {
			lf.applyFilter(line)
		},
	)

	for line := range lf.log.RealTimeLog() {
		lf.applyFilter(line)
	}
}

func (lf viewer) applyFilter(line string) {
	s := strings.Split(line, "|")
	if len(s) != 3 {
		return
	}

	dateMatch, _ := regexp.MatchString(lf.filter.Date, s[0])
	nameMatch, _ := regexp.MatchString(lf.filter.Name, s[2])

	if nameMatch && dateMatch {
		fmt.Println(line)
	}
}

func configureFilter(state FilterState) internal.Filter {
	if filter := state.GetFilter(); filter != nil {
		if proceed := utils.Prompt("Do you want to reuse previous filter? y/n"); proceed == "y" {
			return *filter
		}
	}

	filter := internal.Filter{
		Date: utils.Prompt("Enter date filter (YYYY-MM-DD) or leave blank to match all"),
		Name: utils.Prompt("Enter file name filter (regex) or leave blank to match all"),
	}
	state.SetFilter(filter).Save()

	return filter
}
