package utils

import (
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"log"
	"testing"
)

func TestCEL(t *testing.T) {
	// 环境设定
	env, err := cel.NewEnv(
		//	定义options
		cel.Declarations(
			decls.NewVar("name", decls.String),
			decls.NewVar("group", decls.String)))
	// 分析
	ast, issues := env.Compile(`name.endsWith("1")`)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %s", issues.Err())
	}
	//检查
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("program construction error: %s", err)
	}
	//判断
	out, details, err := prg.Eval(map[string]interface{}{
		"name": "/groups/acme.co/documents/secret-stuff1",
		"group": "acme.co"})
	fmt.Println(out) // 'true'
	fmt.Println(details) // 'true'
}

