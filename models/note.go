package models

import (
	"time"
)

type Note struct {
	ID        int64 
	Text      string            
	Title     string            
	CreatedAt time.Time          
	UpdatedAt time.Time          
	Note_id   string           
}