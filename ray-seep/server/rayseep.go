// @File     : rayseep
// @Author   : Ville
// @Time     : 19-9-26 下午5:25
// server
package server

import (
	"ray-seep/ray-seep/node"
	"ray-seep/ray-seep/server/http"
)

type RaySeepServer struct {
	nodes   map[string]*node.Node
	domains *http.Server
}
