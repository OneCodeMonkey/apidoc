// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"encoding/json"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/docs"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/options"
	"github.com/caixw/apidoc/output/openapi"
)

type marshaler func(v *docs.Docs) ([]byte, error)

// 文档类型定义
const (
	ApidocJSON  = "apidoc+json"
	ApidocYAML  = "apidoc+yaml"
	OpenapiJSON = "openapi+json"
	OpenapiYAML = "openapi+yaml"
	RamlJSON    = "raml+json"
)

var (
	filenames = map[string]string{
		ApidocJSON:  "apidoc.json",
		ApidocYAML:  "apidoc.yaml",
		OpenapiJSON: "openapi.json",
		OpenapiYAML: "openapi.yaml",
		RamlJSON:    "raml.json",
	}
)

// Options 指定了渲染输出的相关设置项。
type Options struct {
	// 文档的保存路径，包含目录和文件名，若为空，则为当前目录下的
	Path string `yaml:"path,omitempty"`

	// 仅输出这些组，为空表示输出所有
	//
	// 若指定的组名实际上不存在，则不会有任何影响。
	Groups []string `yaml:"groups,omitempty"`

	// 输出类型
	Type string `yaml:"type,omitempty"`

	Elapsed time.Duration `yaml:"-"`
	marshal marshaler     // 根据 type 决定转换的函数
}

// Sanitize 对 Options 作一些初始化操作。
func (o *Options) Sanitize() error {
	// TODO 改用默认值
	if o.Path == "" {
		return options.NewFieldError("path", locale.Sprintf(locale.ErrRequired))
	}

	if o.Type == "" {
		o.Type = ApidocJSON
	}

	switch o.Type {
	case ApidocJSON:
		o.marshal = apidocJSONMarshal
	case ApidocYAML:
		o.marshal = apidocYAMLMarshal
	case OpenapiJSON:
		o.marshal = openapi.JSON
	case OpenapiYAML:
		o.marshal = openapi.YAML
	case RamlJSON:
		// TODO
	default:
		return options.NewFieldError("type", locale.Sprintf(locale.ErrInvalidValue))
	}

	return nil
}

// 指定的组名是否包含在输出列表里。
func (o *Options) contains(group string) bool {
	if len(o.Groups) == 0 {
		return true
	}

	for _, g := range o.Groups {
		if g == group {
			return true
		}
	}

	return false
}

func apidocJSONMarshal(v *docs.Docs) ([]byte, error) {
	return json.Marshal(v)
}

func apidocYAMLMarshal(v *docs.Docs) ([]byte, error) {
	return yaml.Marshal(v)
}
