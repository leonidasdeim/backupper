package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/leonidasdeim/backupper/internal"
)

type logFilter struct {
	log internal.Log
	filter
}

type filter struct {
	Date string
	Name string
}

// Runs log filter prompt setup and log viewer runner.
// This function is blocking, should be run from goroutine.
func RunLogFilter(log internal.Log, path string) {
	if proceed := prompt("Do you want to view log? y/n"); proceed != "y" {
		return
	}

	logFilter{
		log:    log,
		filter: configureFilter(path),
	}.runner()
}

func (lf logFilter) runner() {
	internal.Utils.FileScanner(
		lf.log.GetFile(),
		func(line string) {
			lf.applyFilter(line)
		},
	)

	for line := range lf.log.ReadActiveLog() {
		lf.applyFilter(line)
	}
}

func (lf logFilter) applyFilter(line string) {
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

func configureFilter(path string) filter {
	if filter, err := state.load(path); err == nil && filter != nil {
		if proceed := prompt("Do you want to reuse previous filter? y/n"); proceed == "y" {
			return *filter
		}
	}

	filter := filter{
		Date: prompt("Enter date filter (YYYY-MM-DD) or leave blank to match all"),
		Name: prompt("Enter file name filter (regex) or leave blank to match all"),
	}
	state.save(path, filter)

	return filter
}

type _filterState struct{}

var state = _filterState{}

func (_filterState) save(path string, f filter) error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return internal.Utils.OverwriteFile(path, data)
}

func (_filterState) load(path string) (*filter, error) {
	data, err := internal.Utils.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f := filter{}
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}

	return &f, nil
}

func prompt(text string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("- %s: \n", text)
		if answer, _ := reader.ReadString('\n'); answer != "" {
			return strings.TrimSpace(answer)
		}
	}
}
