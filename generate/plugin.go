package generate

import (
	"bytes"
	"fmt"
	"github.com/cloudwego/thriftgo/generator/backend"
	"github.com/cloudwego/thriftgo/generator/golang"
	"github.com/cloudwego/thriftgo/plugin"
	"os"
	"text/template"
	cfg "thrift-gen-docs/config"
)

func newGenerator(req *plugin.Request) (*generator, error) {
	var cfg cfg.Config
	if err := cfg.Unpack(req.PluginParameters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin parameters: %v", err)
	}
	g := &generator{
		buffer:    &bytes.Buffer{},
		config:    &cfg,
		request:   req,
		usedFuncs: make(map[*template.Template]bool),
	}
	lf := backend.LogFunc{
		Info: func(v ...interface{}) {},
		Warn: func(v ...interface{}) {
			g.warnings = append(g.warnings, fmt.Sprint(v...))
		},
		MultiWarn: func(warns []string) {
			g.warnings = append(g.warnings, warns...)
		},
	}
	g.utils = golang.NewCodeUtils(lf)
	g.utils.HandleOptions(req.GeneratorParameters)
	return g, nil
}

func Run(req *plugin.Request) int {
	var warnings []string

	g, err := newGenerator(req)
	if err != nil {
		return exit(&plugin.Response{Warnings: []string{err.Error()}})
	}
	contents, err := g.generate()
	if err != nil {
		return exit(&plugin.Response{Warnings: []string{err.Error()}})
	}

	warnings = append(warnings, g.warnings...)
	for i := range warnings {
		warnings[i] = "[thrift-gen-docs] " + warnings[i]
	}
	res := &plugin.Response{
		Warnings: warnings,
		Contents: contents,
	}

	return exit(res)
}

func exit(res *plugin.Response) int {
	data, err := plugin.MarshalResponse(res)
	if err != nil {
		println("Failed to marshal response:", err.Error())
		return 1
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		println("Error at writing response out:", err.Error())
		return 1
	}
	return 0
}
