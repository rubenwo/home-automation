package script

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/robertkrimen/otto"
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry/models"
	"time"
)

func Log(db *pg.DB) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		m := ""
		for _, v := range call.ArgumentList {
			m = fmt.Sprintf("%s %s", m, v.String())
		}

		_, _ = db.Model(&models.RoutineLog{
			LoggedAt: time.Now(),
			Message:  m,
		}).Insert()

		fmt.Println(m)
		fmt.Println(db)
		// Do what you want with the console output...
		return otto.Value{}
	}
}
