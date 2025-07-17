package controllers

import (
	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
)

func GetAllFoods() ([]models.Food, error) {
	query := `
		SELECT id, name, price, Food)_image, createdat, updatedat, food_id,menu_id 
		FROM foods
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []models.Food
	for rows.Next() {
		var f models.Food
		err := rows.Scan(
			&f.ID,
			&f.Name,
			&f.Price,
			&f.Food_image,
			&f.Created_at,
			&f.Update_at,
			&f.Food_id,
			&f.Menu_id,
		)
		if err != nil {
			return nil, err
		}
		foods = append(foods, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return foods, nil
}
func GetFoodByID(userID string) (*models.Food, error) {
	query := `
		SELECT id, name, price, foodimage,createdat, updatedat, foodid,menuid
		FROM users
		WHERE foodid = ?
	`

	var f models.Food
	err := db.DB.QueryRow(query, userID).Scan(
		&f.ID,
		&f.Name,
		&f.Price,
		&f.Food_image,
		&f.Created_at,
		&f.Update_at,
		&f.Food_id,
		&f.Menu_id,
	)

	if err != nil {
		return nil, err
	}

	return &f, nil
}
