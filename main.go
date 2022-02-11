package main

import (
	DB "courseSelectionSystem/DB"
	"fmt"
)

func main() {
	DB.NewRedisHelper()
	rs := DB.GetRedisHelper()
	res, err := rs.Client.Ping().Result()
	fmt.Println(res, err)
}
