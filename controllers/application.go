package controllers

import (
	"context"
	"os"

	"github.com/StevenSermeus/go-learning/db_access"
	"github.com/StevenSermeus/go-learning/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type ApplicationCreateInput struct {
	Name string `json:"name" binding:"required"`
	RefreshTokenDuration int `json:"refresh_token_duration" binding:"required,min=2,max=2592000"`
	AccessTokenDuration int `json:"access_token_duration" binding:"required,min=1,max=2592000"`
}



func CreateApplication(c *gin.Context) {
	Application := ApplicationCreateInput{}
	err := c.BindJSON(&Application)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if Application.RefreshTokenDuration < Application.AccessTokenDuration {
		c.JSON(400, gin.H{"error": "Refresh token duration should be greater than access token duration"})
		return
	}

	pass_phrase, err := utils.GeneratePassPhrase()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// https://www.rfc-editor.org/rfc/rfc9106.html#name-parameter-choice
	p := &utils.ArgonParams{
		Memory:      2 * 1024,
		Iterations:  1,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	
	hased_pass_phrase, err := utils.GenerateFromPassword(pass_phrase, p)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	api_key, err := utils.GenerateApiKey()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close(ctx)

	queries := db_access.New(conn)
	_, err = queries.CreateApplication(ctx, db_access.CreateApplicationParams{
		Name: Application.Name,
		RefreshTokenDuration: int32(Application.RefreshTokenDuration),
		AccessTokenDuration: int32(Application.AccessTokenDuration),
		PassPhrase: hased_pass_phrase,
		ApiKey: api_key,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"pass_phrase": pass_phrase, "api_key": api_key})
}