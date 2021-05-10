package routines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry/models"
	"io/ioutil"
	"net/http"
)

type RoutineExecutor interface {
	Execute(routine models.Routine) error
}

type TimedRoutineExecutor struct {
}

func (tre *TimedRoutineExecutor) Execute(routine models.Routine) error {
	client := &http.Client{}
	for _, action := range routine.Actions {
		var (
			req *http.Request
			err error
		)
		if action.Data == nil {
			req, err = http.NewRequest(action.Method, action.Addr, nil)
			if err != nil {
				return err
			}
		} else {
			b, err := json.Marshal(action.Data)
			if err != nil {
				return err
			}
			req, err = http.NewRequest(action.Method, action.Addr, bytes.NewBuffer(b))
			if err != nil {
				return err
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		fmt.Println(string(raw))
	}
	return nil
}
