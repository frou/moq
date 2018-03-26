package moq

// moqImports are the imports all moq files get.
var moqImports = []string{}

// moqTemplate is the template for mocked code.
var moqTemplate = `// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package {{.PackageName}}

import (
{{- range .Imports }}
	"{{.}}"
{{- end }}
)

{{ range $i, $obj := .Objects -}}
var (
{{- range .Methods }}
	lock{{$obj.InterfaceName}}Mock{{.Name}}	sync.RWMutex
{{- end }}
)

// {{.InterfaceName}}Mock is a mock implementation of {{.InterfaceName}}.
//
//     func TestSomethingThatUses{{.InterfaceName}}(t *testing.T) {
//
//         // make and configure a mocked {{.InterfaceName}}
//         mocked{{.InterfaceName}} := &{{.InterfaceName}}Mock{ {{ range .Methods }}
//             {{.Name}}Func: func({{ .Arglist }}) {{.ReturnArglist}} {
// 	               panic("TODO: mock out the {{.Name}} method")
//             },{{- end }}
//         }
//
//         // TODO: use mocked{{.InterfaceName}} in code that requires {{.InterfaceName}}
//         //       and then make assertions.
//
//     }
type {{.InterfaceName}}Mock struct {
{{- range .Methods }}
	// {{.Name}}Func mocks the {{.Name}} method.
	{{.Name}}Func func({{ .Arglist }}) {{.ReturnArglist}}
{{ end }}
	// calls tracks calls to the methods.
	calls struct {
{{- range .Methods }}
		// {{ .Name }} holds details about calls to the {{.Name}} method.
		{{ .Name }} []struct {
			{{- range .Params }}
			// {{ .Name | Exported }} is the {{ .Name }} argument value.
			{{ .Name | Exported }} {{ .Type }}
			{{- end }}
		}
{{- end }}
	}
}
{{ range .Methods }}
// {{.Name}} calls {{.Name}}Func.
func (mock *{{$obj.InterfaceName}}Mock) {{.Name}}({{.Arglist}}) {{.ReturnArglist}} {
	if mock.{{.Name}}Func == nil {
		panic("moq: {{$obj.InterfaceName}}Mock.{{.Name}}Func is nil but {{$obj.InterfaceName}}.{{.Name}} was just called")
	}
	callInfo := struct {
		{{- range .Params }}
		{{ .Name | Exported }} {{ .Type }}
		{{- end }}
	}{
		{{- range .Params }}
		{{ .Name | Exported }}: {{ .Name }},
		{{- end }}
	}
	lock{{$obj.InterfaceName}}Mock{{.Name}}.Lock()
	mock.calls.{{.Name}} = append(mock.calls.{{.Name}}, callInfo)
	lock{{$obj.InterfaceName}}Mock{{.Name}}.Unlock()
{{- if .ReturnArglist }}
	return mock.{{.Name}}Func({{.ArgCallList}})
{{- else }}
	mock.{{.Name}}Func({{.ArgCallList}})
{{- end }}
}

// {{.Name}}Calls gets all the calls that were made to {{.Name}}.
// Check the length with:
//     len(mocked{{$obj.InterfaceName}}.{{.Name}}Calls())
func (mock *{{$obj.InterfaceName}}Mock) {{.Name}}Calls() []struct {
		{{- range .Params }}
		{{ .Name | Exported }} {{ .Type }}
		{{- end }}
	} {
	var calls []struct {
		{{- range .Params }}
		{{ .Name | Exported }} {{ .Type }}
		{{- end }}
	}
	lock{{$obj.InterfaceName}}Mock{{.Name}}.RLock()
	calls = mock.calls.{{.Name}}
	lock{{$obj.InterfaceName}}Mock{{.Name}}.RUnlock()
	return calls
}

// {{.Name}}CallCount TODO
func (mock *{{$obj.InterfaceName}}Mock) {{.Name}}CallCount(expectedCount int) error {
	lock{{$obj.InterfaceName}}Mock{{.Name}}.RLock()
	actualCount := len(mock.calls.{{.Name}})
	lock{{$obj.InterfaceName}}Mock{{.Name}}.RUnlock()
	if actualCount == expectedCount {
		return nil
	}

	expectedCountStr := fmt.Sprint(expectedCount)
	if expectedCount == -1 {
		if actualCount >= 1 {
			return nil
		}
		expectedCountStr = "1+"
	}
	return fmt.Errorf(
		"Expected '{{$obj.InterfaceName}}.{{.Name}}' method to be called %s times, but was called %d times",
		expectedCountStr, actualCount)
}
{{ end -}}
{{ end -}}`
