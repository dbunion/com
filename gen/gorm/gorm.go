package gorm

import (
	"bytes"
	"fmt"
	"github.com/dbunion/com/gen"
	"github.com/iancoleman/strcase"
	"github.com/zssky/tc/exec"
	"os"
	"strings"
	"text/template"
)

// GoOrm is gorm model generator
type GoOrm struct {
	cfg gen.Config
}

// NewGoOrm create new gorm with default collection name.
func NewGoOrm() gen.Generator {
	return &GoOrm{}
}

func (g *GoOrm) funcMap() map[string]interface{} {
	return map[string]interface{}{
		"ToLower": strings.ToLower,
		"ToSnake": strcase.ToSnake,
	}
}

// Gen - gen code
func (g *GoOrm) Gen() error {
	if g.cfg.AllInOne {
		return g.genAllInOne()
	}
	return g.genSegregate()
}

func (g *GoOrm) genDB() error {
	file := fmt.Sprintf("gen_%s.go", strcase.ToSnake(g.cfg.Package))
	if g.cfg.GenPath != "" {
		file = fmt.Sprintf("%s/gen_%s.go", g.cfg.GenPath, strcase.ToSnake(g.cfg.Package))
	}
	packageTpl, err := template.New("packageTpl").Funcs(g.funcMap()).Parse(templateDB)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	if err := packageTpl.Execute(&buff, g.cfg); err != nil {
		return err
	}

	writer, err := os.Create(file)
	if err != nil {
		return err
	}

	// write header
	if _, err := fmt.Fprint(writer, buff.String()); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	if _, err := exec.RunShellCommand(fmt.Sprintf("gofmt -w %s", file)); err != nil {
		return err
	}
	return nil
}

func (g *GoOrm) genAllInOne() (gErr error) {
	if err := g.genDB(); err != nil {
		return err
	}

	file := "batch_generated.go"
	if g.cfg.GenPath != "" {
		file = fmt.Sprintf("%s/batch_generated.go", g.cfg.GenPath)
	}
	headerTpl, err := template.New("headerTpl").Funcs(g.funcMap()).Parse(templateHeader)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	if err := headerTpl.Execute(&buff, g.cfg); err != nil {
		return err
	}

	bodyTpl, err := template.New("bodyTpl").Funcs(g.funcMap()).Parse(templateBody)
	if err != nil {
		return err
	}

	writer, err := os.Create(file)
	if err != nil {
		return err
	}

	defer func() {
		if err := writer.Close(); err != nil {
			gErr = err
			return
		}

		if _, err := exec.RunShellCommand(fmt.Sprintf("gofmt -w %s", file)); err != nil {
			gErr = err
			return
		}
	}()

	// write header
	if _, err := fmt.Fprint(writer, buff.String()); err != nil {
		return err
	}

	// gen body
	for i := 0; i < len(g.cfg.Items); i++ {
		buff.Reset()
		item := g.cfg.Items[i]
		if err := bodyTpl.Execute(&buff, item); err != nil {
			return err
		}

		if _, err := fmt.Fprint(writer, buff.String()); err != nil {
			return err
		}
	}

	return nil
}

func (g *GoOrm) genSegregate() (gErr error) {
	if err := g.genDB(); err != nil {
		return err
	}

	for i := 0; i < len(g.cfg.Items); i++ {
		item := g.cfg.Items[i]

		file := fmt.Sprintf("gen_%s.go", strcase.ToSnake(item.Name))
		if g.cfg.GenPath != "" {
			file = fmt.Sprintf("%s/gen_%s.go", g.cfg.GenPath, strcase.ToSnake(item.Name))
		}

		headerTpl, err := template.New("headerTpl").Funcs(g.funcMap()).Parse(templateHeader)
		if err != nil {
			return err
		}

		var buff bytes.Buffer
		if err := headerTpl.Execute(&buff, g.cfg); err != nil {
			return err
		}

		bodyTpl, err := template.New("bodyTpl").Funcs(g.funcMap()).Parse(templateBody)
		if err != nil {
			return err
		}

		writer, err := os.Create(file)
		if err != nil {
			return err
		}

		// write header
		if _, err := fmt.Fprint(writer, buff.String()); err != nil {
			_ = writer.Close()
			return err
		}

		buff.Reset()
		if err := bodyTpl.Execute(&buff, item); err != nil {
			_ = writer.Close()
			return err
		}

		if _, err := fmt.Fprint(writer, buff.String()); err != nil {
			_ = writer.Close()
			return err
		}

		if err := writer.Close(); err != nil {
			return err
		}

		if _, err := exec.RunShellCommand(fmt.Sprintf("gofmt -w %s", file)); err != nil {
			return err
		}
	}

	return nil
}

// StartAndGC start file Helm adapter.
func (g *GoOrm) StartAndGC(config gen.Config) error {
	g.cfg = config
	return nil
}

func init() {
	gen.Register(gen.TypeGormModel, NewGoOrm)
}
