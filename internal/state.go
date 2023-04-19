package internal

import (
	"encoding/json"

	"github.com/leonidasdeim/backupper/internal/utils"
)

type Filter struct {
	Date string `json:"date"`
	Name string `json:"name"`
}

type Directories struct {
	Hot    string `json:"hot"`
	Backup string `json:"backup"`
}

type Data struct {
	Filter      *Filter      `json:"filter"`
	Directories *Directories `json:"dirs"`
}

type State struct {
	path  string
	state Data
}

func LoadState(path string) State {
	p := State{
		path: path,
	}

	data, err := utils.ReadFile(path)
	if err != nil {
		return p
	}

	state := Data{}
	if err := json.Unmarshal(data, &state); err != nil {
		return p
	}

	p.SetState(state)
	return p
}

func (p *State) SetState(state Data) {
	p.state = state
}

func (p *State) SetFilter(filter Filter) *State {
	p.state.Filter = &filter
	return p
}

func (p *State) GetFilter() *Filter {
	return p.state.Filter
}

func (p *State) SetDirectories(dirs Directories) *State {
	p.state.Directories = &dirs
	return p
}

func (p *State) GetDirectories() *Directories {
	return p.state.Directories
}

func (p *State) Save() error {
	data, err := json.Marshal(p.state)
	if err != nil {
		return err
	}
	return utils.OverwriteFile(p.path, data)
}
