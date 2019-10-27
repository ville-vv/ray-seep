package proxy

import (
	"fmt"
	"testing"
)

func TestIdChooser_Choose(t *testing.T) {
	ids := []int64{1245, 12346, 12347, 12348, 12349, 12350}
	cs := &idChooser{}
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
	fmt.Println(cs.Choose(ids))
}
