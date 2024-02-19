package repositories

// import (
// 	"fmt"
// 	"reflect"
// )

// type User struct {
// 	UserID    string `db:"user_id" json:"user_id"`
// 	Email     string `db:"email" json:"email"`
// 	Password  string `db:"password" json:"password"`
// 	FirstName string `db:"first_name" json:"first_name"`
// 	LastName  string `db:"last_name" json:"last_name"`
// 	Active    bool   `db:"active" json:"active"`
// }

// type CreateUserParams struct {
// 	Password  string `db:"password" json:"password"`
// 	FirstName string `db:"first_name" json:"first_name"`
// 	LastName  string `db:"last_name" json:"last_name"`
// 	Active    bool   `db:"active" json:"active"`
// 	Email     string `db:"email" json:"email"`
// }

// func mapCommonFields(model1 interface{}, model2 interface{}) []map[string]interface{} {
// 	result := make([]map[string]interface{}, 0)

// 	val1 := reflect.ValueOf(model1)
// 	val2 := reflect.ValueOf(model2)

// 	typ1 := reflect.TypeOf(model1)
// 	typ2 := reflect.TypeOf(model2)

// 	for i := 0; i < val1.NumField(); i++ {
// 		field1 := typ1.Field(i)
// 		tag1 := field1.Tag.Get("json")

// 		for j := 0; j < val2.NumField(); j++ {
// 			field2 := typ2.Field(j)
// 			tag2 := field2.Tag.Get("json")

// 			if tag1 == tag2 {
// 				data := make(map[string]interface{})
// 				data[tag1] = val1.Field(i).Interface()
// 				data[tag2] = val2.Field(j).Interface()
// 				result = append(result, data)
// 			}
// 		}
// 	}

// 	return result
// }

// func main() {
// 	user := User{
// 		UserID:    "123",
// 		Email:     "user@example.com",
// 		Password:  "password",
// 		FirstName: "John",
// 		LastName:  "Doe",
// 		Active:    true,
// 	}

// 	createUserParams := CreateUserParams{}

// 	result := mapCommonFields(user, createUserParams)
// 	fmt.Println(result)
// }
