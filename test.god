package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

// type User struct {
// 	ID    int    `json:"id"`
// 	Name  string `json:"name"`
// 	Roles []int  `json:"roles"`
// }

// type Person struct {
// 	Name string `json:"name"`
// 	Age  int    `json:"age"`
// }

// func main() {
// 	db := &sql.Sqlite{}
// 	db.DBPath = "./test.db"
// 	err := db.Open(time.Hour)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// }

func main() {
	qqgrouplist := make([]int64, 0, 5)
	qqgrouplist = append(qqgrouplist, 1043728417)
	qqgrouplist = append(qqgrouplist, 1021937014)
	qqgrouplist = append(qqgrouplist, 741433361)
	rand.Shuffle(len(qqgrouplist), func(i, j int) {
		qqgrouplist[i], qqgrouplist[j] = qqgrouplist[j], qqgrouplist[i]
	})
	pigsJSON, _ := json.Marshal(qqgrouplist)
	fmt.Print(string(pigsJSON))

	pigList := make([]int64, 0, 5)
	json.Unmarshal([]byte(string(pigsJSON)), &pigList)
	fmt.Print(pigList)
}
