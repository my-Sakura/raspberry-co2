package reporter

import (
	"bytes"
	"encoding/json"
	"errors"

	"net/http"

	"github.com/dovics/raspberry-co2/util/log"
)

type Reporter interface {
	Register(name string, operate Operate)
	SetTrigger(Trigger) error
}

type Operate func() (interface{}, error)

type Trigger interface {
	Chan() <-chan struct{}
}

func New(url string) Reporter {
	return &defaultReporter{
		target: url,
		client: *http.DefaultClient,
	}
}

type defaultReporter struct {
	operates map[string]Operate
	trigger  Trigger
	running  bool

	target string
	client http.Client
}

func (r *defaultReporter) Run() {
	for range r.trigger.Chan() {
		data := make(map[string]interface{})
		var err error
		for k, operate := range r.operates {
			data[k], err = operate()
			if err != nil {
				data[k] = err
			}
		}
		buf, err := json.Marshal(data)
		if err != nil {
			log.Info(err)
			continue
		}
		resp, err := r.client.Post(r.target, "json", bytes.NewReader(buf))
		if err != nil {
			log.Info(err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Info(errors.New("resp status isn't OK"))
			continue
		}
	}
}

func (r *defaultReporter) SetTrigger(trigger Trigger) error {
	if r.running {
		return errors.New("reporter is running")
	}
	r.trigger = trigger
	return nil
}

func (r *defaultReporter) Register(name string, operate Operate) {
	r.operates[name] = operate
}
