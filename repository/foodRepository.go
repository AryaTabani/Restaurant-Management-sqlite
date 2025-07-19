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
func UpdateFood(food *models.Food) error {
	stmt, err := db.DB.Prepare(`
        UPDATE foods SET name = ?, price = ?, foodimage = ?, menuid = ?, updatedat = ?
        WHERE id = ?
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(food.Name, food.Price, food.Food_image, food.Menu_id, food.Update_at, food.ID)
	return err
}
func GetFoodById(foodid string) (*models.Food, error) {
	var f models.Food
	row := db.DB.QueryRow(`
		SELECT id,name, price, foodimage, createdat, updatedat, foodid, menuid
		FROM foods
		WHERE foodid = ?
	`, foodid)

	err := row.Scan(
		&f.ID, &f.Name, &f.Price, &f.Food_image, &f.Created_at,
		&f.Update_at, &f.Food_id, &f.Menu_id,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
