package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type logFilter struct {
	log Log
	filter
}

type filter struct {
	Date string
	Name string
}

// Runs log filter prompt setup and log viewer runner.
// This function is blocking, should be run from goroutine.
func RunLogFilter(log Log, path string) {
	if proceed := prompt("Do you want to view log? y/n"); proceed != "y" {
		return
	}

	logFilter{
		log:    log,
		filter: resolveFilter(path),
	}.runner()
}

func (lf logFilter) runner() {
	utils.FileScanner(
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

func resolveFilter(path string) filter {
	if f, err := loadState(path); err == nil && f != nil {
		if proceed := prompt("Do you want to reuse previous filter? y/n"); proceed == "y" {
			return *f
		}
	}

	newFilter := filter{
		Date: prompt("Enter date filter (YYYY-MM-DD) or leave blank"),
		Name: prompt("Enter file name filter (regex) or leave blank"),
	}
	saveState(path, newFilter)

	return newFilter
}

func saveState(path string, f filter) error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return utils.OverwriteFile(path, data)
}

func loadState(path string) (*filter, error) {
	data, err := utils.ReadFile(path)
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
