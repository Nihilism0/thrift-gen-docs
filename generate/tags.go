package generate

import (
	"github.com/cloudwego/thriftgo/parser"
	"strings"
)

const (
	Begin = `{
  "swagger": "2.0",
  "info": {
    "title": "swagger",
    "description": "generated by cwgo",`

	parameter1 = `{
            "name": "%s",
            "in": "body",
            "description": "",
            "required": %s,
            "schema": {
              "type": "%s"
            }
          }`

	parameter2 = `{
            "name": "%s",
            "in": "body",
            "description": "",
            "required": %s,
            "schema": {
              "type": "%s",
              "format": "%s"
            }
          }`

	property1 = `"%s": {
                  "type": "%s"
                }`

	property2 = `"%s": {
                  "type": "%s",
                  "format": "%s"
                }`
)

const (
	ApiBaseDomain    = "api.base_domain"
	ApiServiceGroup  = "api.service_group"
	ApiServiceGenDir = "api.service_gen_dir"

	ApiSchemes = "api.scheme"
	ApiConsume = "api.consume"
	ApiProduce = "api.produce"
)

const (
	AnnotationQuery    = "api.query"
	AnnotationForm     = "api.form"
	AnnotationPath     = "api.path"
	AnnotationHeader   = "api.header"
	AnnotationCookie   = "api.cookie"
	AnnotationBody     = "api.body"
	AnnotationRawBody  = "api.raw_body"
	AnnotationJsConv   = "api.js_conv"
	AnnotationNone     = "api.none"
	AnnotationFileName = "api.file_name"

	AnnotationValidator = "api.vd"

	AnnotationGoTag = "go.tag"
)

const (
	ApiGet        = "api.get"
	ApiPost       = "api.post"
	ApiPut        = "api.put"
	ApiPatch      = "api.patch"
	ApiDelete     = "api.delete"
	ApiOptions    = "api.options"
	ApiHEAD       = "api.head"
	ApiAny        = "api.any"
	ApiPath       = "api.path"
	ApiSerializer = "api.serializer"
	ApiGenPath    = "api.handler_path"
	ApiVersion    = "api.version"
)

var (
	HttpMethodAnnotations = map[string]string{
		ApiGet:     "get",
		ApiPost:    "post",
		ApiPut:     "put",
		ApiPatch:   "patch",
		ApiDelete:  "delete",
		ApiOptions: "options",
		ApiHEAD:    "head",
		ApiAny:     "any",
	}

	HttpMethodOptionAnnotations = map[string]string{
		ApiGenPath: "handler_path",
	}

	BindingTags = map[string]string{
		AnnotationPath:    "path",
		AnnotationQuery:   "query",
		AnnotationHeader:  "header",
		AnnotationCookie:  "cookie",
		AnnotationBody:    "json",
		AnnotationForm:    "form",
		AnnotationRawBody: "raw_body",
	}

	SerializerTags = map[string]string{
		ApiSerializer: "serializer",
	}

	ValidatorTags = map[string]string{AnnotationValidator: "vd"}
)

func getAnnotation(input parser.Annotations, target string) []string {
	if len(input) == 0 {
		return nil
	}
	for _, anno := range input {
		if strings.ToLower(anno.Key) == target {
			return anno.Values
		}
		return []string{}
	}
	return []string{}
}

func getAnnotations(input parser.Annotations, targets map[string]string) map[string][]string {
	if len(input) == 0 || len(targets) == 0 {
		return nil
	}
	out := map[string][]string{}
	for k, t := range targets {
		var ret *parser.Annotation
		for _, anno := range input {
			if strings.ToLower(anno.Key) == k {
				ret = anno
				break
			}
		}
		if ret == nil {
			continue
		}
		out[t] = ret.Values
	}
	return out
}
