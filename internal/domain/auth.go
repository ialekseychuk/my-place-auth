package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID         string    
	BusinessID string    
	FirstName  string   
	LastName   string   
	Email      string    	
	Phone      string    	
	Password   string    
	Role       string    // owner, admin, staff
	IsActive   bool     
	CreatedAt  time.Time
	UpdatedAt  time.Time 
}

type JWTClaims struct {
	UserID     string 
	BusinessID string 
	Email      string 
	Role       string 
	IssuedAt   int64  
	ExpiresAt  int64  
	jwt.RegisteredClaims
}