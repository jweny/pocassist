package util

import (
	"container/list"
	"crypto/tls"
	"fmt"
	"github.com/valyala/fasthttp"
	"strings"
	"testing"
	"time"
)


func TestVerifyTargetConnection(t *testing.T) {
	//var m map[string]string
	//
	//if m["aaa"] != "" {
	//	fmt.Println(1)
	//} else {
	//	fmt.Println(0)
	//}

	//originalReq, _ := http.NewRequest("GET", "http://www.jweny.com/", nil)
	//fmt.Println(VerifyTargetConnection(originalReq))


	//req := fasthttp.AcquireRequest()
	//req.SetRequestURI("https://www.360.cn")
	//requestBody := []byte(`{"request":"test"}`)
	//req.SetBody(requestBody)
	//
	//req.Header.SetContentType("application/json")
	//req.Header.SetMethod("POST")
	//
	//url1 := string(req.RequestURI())
	//url2 := req.URI().String()
	//url3 := string(req.Host())
	//url4 := string(req.URI().RequestURI())
	//protocol := string(req.Header.Protocol())
	//url6 := string(req.Header.Header())
	//
	////absRequestURI := strings.HasPrefix(reqURI, "http://") || strings.HasPrefix(reqURI, "https://")
	//fmt.Println(url1)
	//fmt.Println(url2)
	//fmt.Println(url3)
	//fmt.Println(url4)
	//fmt.Println(protocol)
	//fmt.Println(url6)
	/*
	https://www.360.cn
	https://www.360.cn/
	www.360.cn
	/
	*/

}

func TestCopyRequest(t *testing.T) {
	test := "ccc"
	str := fmt.Sprintf("%s%s%s", test, "/", strings.TrimPrefix("/aaa", "/"))
	curPath := fmt.Sprint(test, "/" ,strings.TrimPrefix("/aaa", "/"))
	fmt.Println(str)
	fmt.Println(curPath)
}

var Default = -1024

type Node struct {
	data interface{}
	lchild *Node
	rchild *Node
}

func InsertNodeToTree(tree *Node, node *Node){
	if tree == nil {
		return
	}
//		root节点
	if tree.data == Default{
		tree.data = node.data
		return
	}
	if node.data.(int) > tree.data.(int) {
		if tree.lchild == nil {
			tree.lchild = &Node{data : Default}
		}
		InsertNodeToTree(tree.lchild, node)
	}
	if node.data.(int) < tree.data.(int) {
		if tree.rchild == nil {
			tree.rchild = &Node{data : Default}
		}
		InsertNodeToTree(tree.rchild, node)
	}
}

func InitTree(values ...int) *Node{
	root := Node{data:Default,lchild:nil,rchild:nil}
	for _, d := range values {
		n := Node{data:d}
		InsertNodeToTree(&root,&n)
	}
	return &root
}

func PreOrderTraverse(node *Node){
	if node == nil {
		return
	}
	fmt.Println(node.data)
	PreOrderTraverse(node.lchild)
	PreOrderTraverse(node.rchild)
}

//层序遍历
func LevelOrderTraverse(node *Node) {
	if node == nil {
		return
	}
	query := list.New()
	query.PushBack(node)
	if query.Len() >0 {
		// 队首出列
		head := query.Remove(query.Front())
		node := head.(*Node)
		fmt.Println(node.data)
		if node.lchild != nil {
			query.PushBack(node.lchild)
		}
		if node.lchild != nil {
			query.PushBack(node.rchild)
		}
	}
}

func TestDealMultipart(t *testing.T) {
	treeNode := InitTree(5, 4, 6, 8, 9, 7, 1, 3, 2)
	fmt.Println(treeNode)
	fmt.Println("______________")
	LevelOrderTraverse(treeNode)
	fmt.Println("______________")
	PreOrderTraverse(treeNode)
	fmt.Println("______________")
}

func TestDoFasthttpRequest(t *testing.T) {
	client := &fasthttp.Client{
		// If InsecureSkipVerify is true, TLS accepts any certificate
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
	}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI("http://47.104.152.145:29999/api/settings/values")
	client.DoTimeout(req, resp, time.Duration(5)*time.Second)
}