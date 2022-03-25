package controllers

import (
	"strconv"

	"github.com/revel/revel"
)

type UserController struct {
	*revel.Controller
}

func GetAllUsers(c UserController) revel.Result {
	db := connect()
	defer db.Close()

	query := "SELECT * from users"

	rows, err := db.Query(query)
	if err != nil {
		return nil
	}

	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address); err != nil {
			return nil
		} else {
			users = append(users, user)
		}
	}
	if len(users) > 1 {
		var response UsersResponse
		response.Status = 200
		response.Message = "Sucess"
		response.Data = users
		return c.RenderJSON(response)
	} else if len(users) == 1 {
		var response UserResponse
		response.Status = 200
		response.Message = "Sucess"
		response.Data = user
		return c.RenderJSON(response)
	} else {
		var response UserResponse
		response.Status = 400
		response.Message = "No User Found"
		return c.RenderJSON(response)
	}
}
func InsertNewUser(c UserController) revel.Result {
	db := connect()
	defer db.Close()
	var user User
	user.Name = c.Params.Form.Get("name")
	user.Age, _ = strconv.Atoi(c.Params.Form.Get("age"))
	user.Address = c.Params.Form.Get("address")

	result, _ := db.Exec("insert into users (name, age, address) values (?, ?, ?)", user.Name, user.Age, user.Address)

	num, _ := result.RowsAffected()

	var response UserResponse
	if num != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = user
		c.Request.Header.Set("Content-Type", "application/json")
		return c.RenderJSON(response)
	} else {
		response.Status = 400
		response.Message = "Bad Request"
		return c.RenderJSON(response)
	}
}
func DeleteUser(c UserController) revel.Result {
	db := connect()
	defer db.Close()

	userID := c.Params.Route.Get("id")

	_, errQuery := db.Exec("DELETE FROM users WHERE id=?", userID)

	if errQuery == nil {
		var response UserResponse
		response.Status = 400
		response.Message = "Delete Failed"
		return c.RenderJSON(response)
	} else {
		var response UserResponse
		response.Status = 200
		response.Message = "Delete Sucess"
		return c.RenderJSON(response)
	}
}
func UpdateUser(c UserController) revel.Result {
	{
		db := connect()
		defer db.Close()

		var user User = GetUser(c.Params.Route.Get("id"))
		if (c.Params.Form.Get("name")) != "" {
			user.Name = c.Params.Form.Get("name")
		}
		if c.Params.Form.Get("address") != "" {
			user.Address = c.Params.Form.Get("address")
		}
		var temp, _ = strconv.Atoi(c.Params.Form.Get("age"))
		if temp != 0 {
			user.Age, _ = strconv.Atoi(c.Params.Form.Get("age"))
		}

		result, _ := db.Exec("UPDATE users SET name = ?, age = ?, address = ? WHERE id = ?", user.Name, user.Age, user.Address, user.ID)

		num, _ := result.RowsAffected()

		var response UserResponse
		if num != 0 {
			response.Status = 200
			response.Message = "Success"
			response.Data = user
			return c.RenderJSON(response)
		} else {
			response.Status = 200
			response.Message = "Success"
			response.Data = user
			return c.RenderJSON(response)
		}
	}
}

func GetUser(user_id string) User {
	db := connect()
	defer db.Close()
	var user User
	query := "SELECT * from users WHERE ID = " + user_id
	rows, _ := db.Query(query)
	if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address); err != nil {
		return user
	}
	return user
}
