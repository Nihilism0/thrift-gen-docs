package generate

import (
	"bytes"
	"fmt"
	"github.com/cloudwego/hertz/cmd/hz/util"
	"github.com/cloudwego/thriftgo/generator/golang"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/plugin"
	"path/filepath"
	"text/template"
	"thrift-gen-docs/config"
)

type generator struct {
	config     *config.Config
	request    *plugin.Request
	utils      *golang.CodeUtils
	buffer     *bytes.Buffer
	indentNum  int
	warnings   []string
	usedFuncs  map[*template.Template]bool
	enumImport []string
}

func (g *generator) generate() ([]*plugin.Generated, error) {
	var ret []*plugin.Generated
	asts := g.request.AST.DepthFirstSearch()
	for ast := range asts {
		g.buffer.Reset()
		g.enumImport = g.enumImport[:0]
		scope, err := golang.BuildScope(g.utils, ast)
		if err != nil {
			return nil, err
		}
		g.writeLine(Begin)
		ss := ast.GetServices()
		for _, s := range ss {
			g.indent()
			g.indent()
			apiVersions := getAnnotation(s.Annotations, ApiVersion)
			if len(apiVersions) > 1 {
				return nil, fmt.Errorf("too many 'api.XXX' annotations: %s", apiVersions)
			} else if len(apiVersions) == 1 {
				g.writeLinef("\"version\": \"%s\"\n", apiVersions[0])
				g.unindent()
				g.writeLine("},")
			} else {
				g.writeLine("\"version\": \"\"")
				g.unindent()
				g.writeLine("},")
			}
			domainAnno := getAnnotation(s.Annotations, ApiBaseDomain)
			if len(domainAnno) == 1 {
				g.writeLinef("\"host\":\"%s\",\n", domainAnno[0])
			} else {
				g.writeLine("\"host\":\"\",")
			}
			schemesAnno := getAnnotation(s.Annotations, ApiSchemes)
			g.writeLine("\"schemes\": [")
			g.indent()
			if len(schemesAnno) >= 1 {
				for k, scheme := range schemesAnno {
					g.writeLinef("\"%s\"\n", scheme[k])
				}
			}
			g.unindent()
			g.writeLine("],")
			g.writeLine("\"paths\": {")
			g.indent()
			fs := s.GetFunctions()
			//methods := make([]*gen.HttpMethod, 0, len(ms))
			//clientMethods := make([]*gen.ClientMethod, 0, len(ms))
			for _, f := range fs {
				//add the global sets
				rs := getAnnotations(f.Annotations, HttpMethodAnnotations)
				if len(rs) > 1 {
					return nil, fmt.Errorf("too many 'api.XXX' annotations: %s", rs)
				}
				if len(rs) == 0 {
					continue
				}
				hmethod, path := util.GetFirstKV(rs)
				if len(path) != 1 || path[0] == "" {
					return nil, fmt.Errorf("invalid api.%s  for %s.%s: %s", hmethod, s.Name, f.Name, path)
				}
				g.writeLinef("\"%s\":{\n", path[0])
				g.indent()
				g.writeLinef("\"%s\": {\n", hmethod)
				g.indent()
				//add consumes or produces
				contents := getAnnotation(f.Annotations, ApiConsume)
				if len(contents) == 0 {
					contents = getAnnotation(f.Annotations, ApiProduce)
					if len(contents) >= 1 {
						g.writeLine("\"produces\": [")
						g.indent()
						for k, content := range contents {
							if k != 0 {
								g.writeLine(",")
							}
							g.writeLinef("\"%s\"", content)
						}
						g.writeLine("")
						g.unindent()
						g.writeLine("],")
					}
				} else {
					g.writeLine("\"consumes\": [")
					g.indent()
					for k, content := range contents {
						if k != 0 {
							g.writeLine(",")
						}
						g.writeLinef("\"%s\"", content)
					}
					g.writeLine("")
					g.unindent()
					g.writeLine("],")
				}
				//add parameters
				g.writeLine("\"parameters\": [")
				g.indent()
				name := f.Arguments[0].Type.Name
				names := ast.Structs
				re := findStructByName(names, name)
				if re == nil {
					return nil, fmt.Errorf("invalid request struct : %s", name)
				}
				for k, filed := range re.Fields {
					if k != 0 {
						g.buffer.WriteString(",\n")
					}
					var isRequired string
					if filed.Requiredness.IsRequired() {
						isRequired = "true"
					} else {
						isRequired = "false"
					}
					parameter, err := GetParameter(filed.Name, isRequired, filed.Type.Name)
					if err != nil {
						return nil, err
					}
					g.writeLinef(parameter)
				}
				g.unindent()
				g.writeLine("")
				g.writeLine("],")
				//add response
				g.writeLine("\"responses\": {")
				g.indent()
				g.writeLine("\"200\": {")
				g.indent()
				g.writeLine("\"description\": \"\",")
				g.writeLine("\"schema\": {")
				g.indent()
				g.writeLine("\"type\": \"object\",")
				g.writeLine("\"properties\": {")
				g.indent()
				name = f.FunctionType.Name
				re = findStructByName(names, name)
				if re == nil {
					return nil, fmt.Errorf("invalid response struct : %s", name)
				}
				for k, filed := range re.Fields {
					if k != 0 {
						g.buffer.WriteString(",\n")
					}
					property, err := GetProperty(filed.Name, filed.Type.Name)
					if err != nil {
						return nil, err
					}
					g.writeLinef(property)
				}
			}
			g.unindent()
			g.writeLine("")
			g.writeLine("}")
			g.unindent()
			g.writeLine("}")
			g.unindent()
			g.writeLine("}")
			g.unindent()
			g.writeLine("}")
			g.unindent()
			g.writeLine("}")
			g.unindent()
			g.writeLine("}")
			g.unindent()
			g.writeLine("}")
			g.unindent()
			g.writeLine("}")

		}
		if err != nil {
			return nil, err
		}
		g.utils.SetRootScope(scope)
		fp := "swagger.json"
		full := filepath.Join(g.request.OutputPath, fp)
		ret = append(ret, &plugin.Generated{
			Content: g.buffer.String(),
			Name:    &full,
		})
	}
	return ret, nil
}

func (g *generator) writeLine(str string) {
	for i := 0; i < g.indentNum; i++ {
		g.buffer.WriteString("\t")
	}
	g.buffer.WriteString(str + "\n")
}

func (g *generator) writeLinef(format string, args ...interface{}) {
	for i := 0; i < g.indentNum; i++ {
		g.buffer.WriteString("\t")
	}
	g.buffer.WriteString(fmt.Sprintf(format, args...))
}

func (g *generator) indent() {
	g.indentNum++
}

func (g *generator) unindent() {
	g.indentNum--
}

func findStructByName(structs []*parser.StructLike, searchName string) *parser.StructLike {
	for _, s := range structs {
		if s.Name == searchName {
			return s
		}
	}
	return nil
}

func GetParameter(name, required, filedName string) (string, error) {
	if filedName == "i16" || filedName == "i32" || filedName == "i64" || filedName == "byte" {
		if filedName == "i16" {
			filedName = "int16"
		} else if filedName == "i32" {
			filedName = "int32"
		} else if filedName == "i64" {
			filedName = "int64"
		} else {
			filedName = "int8"
		}
		return formatStr(parameter2, name, required, "integer", filedName), nil
	} else if filedName == "bool" {
		return formatStr(parameter1, name, required, "boolean"), nil
	} else if filedName == "double" {
		return formatStr(parameter2, name, required, "number", "double"), nil
	} else if filedName == "string" {
		return formatStr(parameter1, name, required, "string"), nil
	} else {
		return "", fmt.Errorf("invalid parameter : %s ", filedName)
	}
}

func GetProperty(name, filedName string) (string, error) {
	if filedName == "i16" || filedName == "i32" || filedName == "i64" || filedName == "byte" {
		if filedName == "i16" {
			filedName = "int16"
		} else if filedName == "i32" {
			filedName = "int32"
		} else if filedName == "i64" {
			filedName = "int64"
		} else {
			filedName = "int8"
		}
		return formatStr(property2, name, "integer", filedName), nil
	} else if filedName == "bool" {
		return formatStr(property1, name, "boolean"), nil
	} else if filedName == "double" {
		return formatStr(property2, name, "number", "double"), nil
	} else if filedName == "string" {
		return formatStr(property1, name, "string"), nil
	} else {
		return "", fmt.Errorf("invalid property : %s ", filedName)
	}
}

func formatStr(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
