package exporter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Operate func() (interface{}, error)
type Exporter interface {
	Run() error
	Register(name string, operate Operate)
}

type defaultExporter struct {
	operates map[string]Operate
}

func NewExporter() Exporter {
	return &defaultExporter{
		operates: make(map[string]Operate),
	}
}

func (e *defaultExporter) Register(name string, operate Operate) {
	e.operates[name] = operate
}

func (e *defaultExporter) Run() error {
	handlerWrapper := func(operate Operate) gin.HandlerFunc {
		return func(c *gin.Context) {
			data, err := operate()
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"error": err})
				return 
			}

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": data})
		}
	}

	engine := gin.Default()
	for k, o := range e.operates {
		engine.GET(k, handlerWrapper(o))
	}

	return engine.Run()
}
