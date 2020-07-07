package main

import (
	"fmt"
	"ray-seep/run/admin/input"
	"ray-seep/run/admin/store"
)

func main() {
	db := store.NewMysqlStore(input.InputStoreInfo())
	for {
		ac, pt := input.UserAccountInfo()
		if err := db.AddRaySeepUser(ac, pt); err != nil {
			fmt.Println("添加账户错误：", err)
			return
		}
	}
}
