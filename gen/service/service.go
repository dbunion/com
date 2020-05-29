package service

import (
	"bytes"
	"fmt"
	"github.com/dbunion/com/gen"
	"github.com/iancoleman/strcase"
	"github.com/zssky/tc/exec"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

// Bot service generator
type Bot struct {
	cfg gen.Config
}

// NewServiceBot create new service bot with default collection name.
func NewServiceBot() gen.Generator {
	return &Bot{}
}

func (b *Bot) funcMap() map[string]interface{} {
	return map[string]interface{}{
		"ToLower": strings.ToLower,
		"ToSnake": strcase.ToSnake,
		"MakePreload": func(relations []string) string {
			sql := ""
			for i := 0; i < len(relations); i++ {
				sql += fmt.Sprintf("Preload(\"%s\").", relations[i])
			}
			sql = strings.TrimSuffix(sql, ".")
			return sql
		},
	}
}

// Gen - gen code
func (b *Bot) Gen() error {
	if b.cfg.AllInOne {
		return b.genAllInOne()
	}
	if err := b.genCommon(); err != nil {
		return err
	}

	return b.genSegregate()
}

func (b *Bot) genCommon() error {
	file := fmt.Sprintf("gen_%s.go", strcase.ToSnake(b.cfg.Package))
	if b.cfg.GenPath != "" {
		file = fmt.Sprintf("%s/gen_%s.go", b.cfg.GenPath, strcase.ToSnake(b.cfg.Package))
	}

	commonTpl, err := template.New("serviceCommonTpl").Funcs(b.funcMap()).Parse(serviceCommon)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	if err := commonTpl.Execute(&buff, b.cfg); err != nil {
		return err
	}

	writer, err := os.Create(file)
	if err != nil {
		return err
	}

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

func (b *Bot) genAllInOne() error {
	return fmt.Errorf("not support all in one")
}

func (b *Bot) genSegregate() (gErr error) {
	for i := 0; i < len(b.cfg.ServiceCfg.Items); i++ {
		item := b.cfg.ServiceCfg.Items[i]

		reqType := reflect.TypeOf(item.Req)
		dstType := reflect.TypeOf(item.Dst)

		reqName, dstName := reqType.Name(), dstType.Name()
		if reqType.Kind() == reflect.Ptr {
			reqName = reqType.Elem().Name()
		}
		if dstType.Kind() == reflect.Ptr {
			dstName = dstType.Elem().Name()
		}

		req := item.Req
		dst := item.Dst
		index := item.Index

		file := fmt.Sprintf("gen_%s.go", strcase.ToSnake(dstName))
		if b.cfg.GenPath != "" {
			file = fmt.Sprintf("%s/gen_%s.go", b.cfg.GenPath, strcase.ToSnake(dstName))
		}

		var buff bytes.Buffer
		bodyTpl, err := template.New("bodyTpl").Funcs(b.funcMap()).Parse(serviceTemplate)
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

		if err := bodyTpl.Execute(&buff, map[string]interface{}{
			"ImportModelPath":  b.cfg.ServiceCfg.ImportModelPath,
			"Package":          b.cfg.Package,
			"DstName":          dstName,
			"ReqName":          reqName,
			"Req":              req,
			"Dst":              dst,
			"DstConstruct":     makeConstruct(req, dst),
			"ConvertConstruct": makeConstruct(dst, req),
			"Index":            index,
			"CheckApp":         item.CheckApp,
		}); err != nil {
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

func constructField(reqField reflect.StructField, dstField reflect.StructField) string {
	needCheckType := false
	if strings.EqualFold(strings.ToLower(reqField.Name), strings.ToLower(dstField.Name)) {
		if reqField.Type.Kind() == dstField.Type.Kind() {
			return fmt.Sprintf("%s: req.%s,\n", dstField.Name, reqField.Name)
		}
		needCheckType = true
	}

	if needCheckType {
		switch dstField.Type.Kind() {
		case reflect.Uint:
			return fmt.Sprintf("%s: uint(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Uint8:
			return fmt.Sprintf("%s: uint8(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Uint16:
			return fmt.Sprintf("%s: uint16(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Uint32:
			return fmt.Sprintf("%s: uint32(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Uint64:
			return fmt.Sprintf("%s: uint64(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Int:
			return fmt.Sprintf("%s: int(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Int8:
			return fmt.Sprintf("%s: int8(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Int16:
			return fmt.Sprintf("%s: int16(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Int32:
			return fmt.Sprintf("%s: int32(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Int64:
			return fmt.Sprintf("%s: int64(req.%s),\n", dstField.Name, reqField.Name)
		case reflect.Struct:

		default:
			return fmt.Sprintf("%s: req.%s,\n", dstField.Name, reqField.Name)
		}
	}
	return ""
}

func makeConstruct(req interface{}, dst interface{}) string {
	reqTyp := reflect.TypeOf(req)
	dstTyp := reflect.TypeOf(dst)

	if reqTyp.Kind() == reflect.Ptr {
		reqTyp = reqTyp.Elem()
	}

	if dstTyp.Kind() == reflect.Ptr {
		dstTyp = dstTyp.Elem()
	}

	var buffer bytes.Buffer

	// write struct
	dstBase := filepath.Base(dstTyp.PkgPath())
	fmt.Fprintf(&buffer, "val := &%s.%s{\n", dstBase, dstTyp.Name())
	for i := 0; i < reqTyp.NumField(); i++ {
		for j := 0; j < dstTyp.NumField(); j++ {
			code := constructField(reqTyp.Field(i), dstTyp.Field(j))
			if code != "" {
				fmt.Fprint(&buffer, code)
			}
		}
	}
	fmt.Fprintf(&buffer, "}\n")

	return buffer.String()
}

// StartAndGC start file Helm adapter.
func (b *Bot) StartAndGC(config gen.Config) error {
	b.cfg = config
	return nil
}

func init() {
	gen.Register(gen.TypeService, NewServiceBot)
}
