package service

import (
	"github.com/dbunion/com/gen"
	"testing"
)

// ServiceApp - service app
type ServiceApp struct {
	ID int64 `json:"id,omitempty"`
	// @inject_tag: binding:"required"`
	Name string `json:"name,omitempty"`
	// @inject_tag: binding:"required"`
	Alias      string `json:"alias,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
	UpdateTime string `json:"update_time,omitempty"`
}

// ModelApp - model app info
type ModelApp struct {
	Name   string
	Alias  string
	Remark string
}

func TestGenSegregate(t *testing.T) {
	appItem := gen.SItem{
		Req:   &ServiceApp{},
		Dst:   &ModelApp{},
		Index: 15000,
		CheckApp: true,
	}

	generator, err := gen.NewGenerator(gen.TypeService, gen.Config{
		Package:  "service",
		GenPath:  "/tmp",
		AllInOne: false,
		ServiceCfg: gen.ServiceGenConfig{
			ImportPaths: []string{
				"github.com/dbunion/proto",
				"github.com/dbunion/com",
			},
			Items: []gen.SItem{
				appItem,
			},
		},
	})

	if err != nil {
		t.Fatalf("create new generator error, err:%v", err)
	}

	if err := generator.Gen(); err != nil {
		t.Fatalf("gen code error, err:%v", err)
	}
}
