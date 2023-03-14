package web_copy

import (
	"strings"
)

type router struct {
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: make(map[string]*node),
	}
}

func (r *router) addRouter(method string, path string, handlerFunc HandlerFunc) {
	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	if path[0] != '/' && path[len(path)-1] == '/' {
		panic("web：路由错误")
	}

	root, ok := r.trees[method]
	if !ok {
		r.trees[method] = &node{
			path: "/",
		}
		root = r.trees[method]
	}

	if path == "/" {
		root.handlerFunc = handlerFunc
		return
	}

	ages := strings.Split(path[1:], "/")
	for _, v := range ages {
		if root.childNode == nil {
			root.childNode = make(map[string]*node)
		}
		tem := root.childOrCreate(v)
		root.childNode[v] = tem
		root = tem
	}

	root.handlerFunc = handlerFunc
}

func (r *router) findRouter(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{n: root}, true
	}

	ages := strings.Split(strings.Trim(path, "/"), "/")
	mi := &matchInfo{}
	for _, v := range ages {
		var paramBool bool
		root, paramBool, ok = root.childOf(v)
		if !ok {
			return nil, false
		}
		if paramBool {
			mi.addValue(root.path[1:], v)
		}
	}

	mi.n = root
	return mi, true
}

func (n *node) childOrCreate(path string) *node {
	if path == "" {
		panic("web：非法路由")
	}

	if path == "*" {
		if n.paramPath != nil {
			panic("web：不允许同时注册通配符路由和参数路由")
		}
		if n.starPath == nil {
			n.starPath = &node{
				path: "*",
			}
		}
		return n.starPath
	}

	if path[0] == ':' {
		if n.starPath != nil {
			panic("web：不允许同时注册通配符路由和参数路由")
		}
		if n.paramPath != nil {
			if n.paramPath.path != path {
				panic("web：已经注册过该参数路由了")
			}
		} else {
			n.paramPath = &node{
				path: path,
			}
		}
		return n.paramPath
	}

	next, ok := n.childNode[path]
	if !ok {
		next = &node{
			path: path,
		}
	}

	return next
}

func (n *node) childOf(path string) (*node, bool, bool) {
	if path == "" {
		return nil, false, false
	}

	if n.childNode == nil {
		if n.paramPath != nil {
			return n.paramPath, true, true
		}
		return n.starPath, false, n.starPath == nil
	}
	child, ok := n.childNode[path]
	if !ok {
		if n.paramPath != nil {
			return n.paramPath, true, true
		}
		return n.starPath, false, n.starPath == nil
	}
	return child, false, true
}

type node struct {
	path        string
	handlerFunc HandlerFunc
	childNode   map[string]*node

	// 通配符
	starPath *node
	// 参数路由
	paramPath *node
}

type matchInfo struct {
	n         *node
	paramPath map[string]string
}

func (m *matchInfo) addValue(key, value string) {
	if m.paramPath == nil {
		m.paramPath = make(map[string]string)
	}
	m.paramPath[key] = value
}
