package repository

import (
	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
)

func CreateFood(food *models.Food) (int64, error) {
	stmt, err := db.DB.Prepare(`
		INSERT INTO foods (name, price, foodimage, createdat, updatedat, foodid, menuid)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		food.Name, food.Price, food.Food_image,
		food.Created_at, food.Update_at, food.Food_id, food.Menu_id,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
