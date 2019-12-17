package databus

import (
	"fmt"
	"ray-seep/ray-seep/server/env_init"
	"testing"
)

func TestSnakeString(t *testing.T) {
	src := "RaySeepID"
	fmt.Println(env_init.SnakeToCameString(src))
}
func TestCamelString(t *testing.T) {
	fmt.Println(env_init.CamelToSnakeString("ray_seep_id"))
}
