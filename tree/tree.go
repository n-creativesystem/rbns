package tree

import (
	"regexp"
	"strings"
	"sync"
)

type Node struct {
	path        string
	permissions []string
	children    map[string]*Node
}

type Tree struct {
	nodes map[string]*Node
}

const (
	pathRoot           = "/"
	pathDelimiter      = "/"
	parameterDelimiter = ":"
	ptnLeftDelimiter   = "["
	ptnRightDelimiter  = "]"
	ptnWildCard        = "(.+)"
)

func NewTree() *Tree {
	return &Tree{
		nodes: map[string]*Node{},
	}
}

func (t *Tree) Insert(method string, url string, permissions ...string) error {
	current, ok := t.nodes[method]
	if !ok {
		t.nodes[method] = &Node{
			children: map[string]*Node{},
		}
		current = t.nodes[method]
	}
	if url == pathRoot {
		current.path = url
		current.permissions = make([]string, 0, len(permissions))
		for _, p := range permissions {
			if p != "" {
				current.permissions = append(current.permissions, p)
			}
		}
		return nil
	}
	paths := parseUrl(url)
	for i, path := range paths {
		nextNode, ok := current.children[path]
		if ok {
			current = nextNode
		} else {
			current.children[path] = &Node{
				path:        path,
				permissions: []string{},
				children:    map[string]*Node{},
			}
			current = current.children[path]
		}

		if i == len(paths)-1 {
			current.path = path
			current.permissions = make([]string, 0, len(permissions))
			for _, p := range permissions {
				if p != "" {
					current.permissions = append(current.permissions, p)
				}
			}
			break
		}
	}
	return nil
}

func (t *Tree) Search(method string, url string) []string {
	current, ok := t.nodes[method]
	if !ok {
		return nil
	}
	paths := parseUrl(url)
	for _, path := range paths {
		nextNode, ok := current.children[path]
		if ok {
			current = nextNode
			continue
		}

		// 子ノードなし
		if len(current.children) == 0 {
			if current.path != path {
				return nil
			}
			break
		}

		isParamMatch := false
		for path := range current.children {
			if string([]rune(path)[0]) == parameterDelimiter {
				ptn := getPattern(path)
				reg, err := regCache.get(ptn)
				if err != nil {
					return nil
				}
				if reg.Match([]byte(path)) {
					current = current.children[path]
					isParamMatch = true
					break
				}
				return nil
			}
		}

		if !isParamMatch {
			// 一致するパスがない
			return nil
		}
	}
	permissions := make([]string, 0, len(current.permissions))
	permissions = append(permissions, current.permissions...)
	return permissions
}

func parseUrl(url string) []string {
	paths := strings.Split(url, pathDelimiter)
	res := make([]string, 0, len(paths))
	for _, p := range paths {
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}

// :id[^\d+$] -> ^\d+$
// :id -> (.+)
func getPattern(path string) string {
	leftIdx := strings.Index(path, ptnLeftDelimiter)
	rightIdx := strings.Index(path, ptnRightDelimiter)

	if leftIdx == -1 || rightIdx == -1 {
		return ptnWildCard
	}

	return path[leftIdx+1 : rightIdx]
}

var regCache = newRegexpCache()

func newRegexpCache() *regexpCache {
	return &regexpCache{
		value: make(map[string]*regexp.Regexp),
	}
}

type regexpCache struct {
	mu    sync.Mutex
	value map[string]*regexp.Regexp
}

func (rc *regexpCache) get(ptn string) (*regexp.Regexp, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	v, ok := rc.value[ptn]
	if ok {
		return v, nil
	}
	reg, err := regexp.Compile(ptn)
	if err != nil {
		return nil, err
	}
	rc.value[ptn] = reg
	return reg, nil
}
