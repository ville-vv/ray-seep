package dao

import (
	"fmt"
	"testing"
)

func TestSnakeString(t *testing.T) {
	src := "RaySeepID"
	fmt.Println(SnakeToCameString(src))
}
func TestCamelString(t *testing.T) {
	fmt.Println(CamelToSnakeString("ray_seep_id"))
}
