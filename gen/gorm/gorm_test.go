package gorm

import (
	"github.com/dbunion/com/gen"
	"testing"
)

func TestGenAllInOne(t *testing.T) {
	generator, err := gen.NewGenerator(gen.TypeGormModel, gen.Config{
		Package:  "gorm",
		GenPath:  "/tmp",
		AllInOne: true,
		ModelCfg: gen.ModelGenConfig{
			Items: []gen.Item{
				{
					Name: "UserAgent",
					Detail: `type UserAgent struct {
	Model
	Name string
	Detail string
}`,
				},
				{
					Name:      "TestCase",
					Relations: []string{"UserAgent"},
					Detail: `type TestCase struct {
	Model
	Name string
	Range int64
}
`,
				},
			},
			MaxIdleConns: 5,
			MaxOpenConns: 10,
		},
	})

	if err != nil {
		t.Fatalf("create new generator error, err:%v", err)
	}

	if err := generator.Gen(); err != nil {
		t.Fatalf("gen code error, err:%v", err)
	}

}

func TestGenSegregate(t *testing.T) {
	generator, err := gen.NewGenerator(gen.TypeGormModel, gen.Config{
		Package:  "gorm",
		GenPath:  "/tmp",
		AllInOne: false,
		ModelCfg: gen.ModelGenConfig{
			Items: []gen.Item{
				{
					Name: "UserAgent",
					Detail: `type UserAgent struct {
	Model
	Name string
	Detail string
}`,
				},
				{
					Name:      "TestCase",
					Relations: []string{"UserAgent", "TestCase"},
					Detail: `type TestCase struct {
	Model
	Name string
	Range int64
}
`}},
			MaxIdleConns: 5,
			MaxOpenConns: 10,
		},
	})

	if err != nil {
		t.Fatalf("create new generator error, err:%v", err)
	}

	if err := generator.Gen(); err != nil {
		t.Fatalf("gen code error, err:%v", err)
	}
}
