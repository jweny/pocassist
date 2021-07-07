package scripts

import (
	"fmt"
	log "github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"strings"
)

type ScriptScanArgs struct {
	Host    string
	Port    uint16
	IsHTTPS bool
}

type ScriptScanFunc func(args *ScriptScanArgs) (*util.ScanResult, error)

var scriptHandlers = map[string]ScriptScanFunc{}

// GetScriptFunc 返回 pocName 对应的方法
func GetScriptFunc(pocName string) ScriptScanFunc {
	if f, ok := scriptHandlers[pocName]; ok {
		return f
	}
	return nil
}

func ScriptRegister(pocName string, handler ScriptScanFunc) {
	if _, ok := scriptHandlers[pocName]; ok {
		log.Panic("[script register vulId ]", pocName)
	}
	scriptHandlers[pocName] = handler
}

func ConstructUrl(args *ScriptScanArgs, uri string) string {
	var rawUrl string
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}
	var scheme string
	if args.IsHTTPS {
		scheme = "https"
	} else {
		scheme = "http"
	}
	if args.Port == 80 || args.Port == 443 {
		rawUrl = fmt.Sprintf("%v://%v%v", scheme, args.Host, uri)
	} else {
		rawUrl = fmt.Sprintf("%v://%v:%v%v", scheme, args.Host, args.Port, uri)
	}
	return rawUrl
}




