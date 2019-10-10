// @File     : rayseep
// @Author   : Ville
// @Time     : 19-9-26 下午5:25
// server
package server

import (
	"ray-seep/ray-seep/server/http"
	"ray-seep/ray-seep/server/node"
)

type RaySeepServer struct {
	nodes   map[string]*node.Node
	domains *http.Server
}
