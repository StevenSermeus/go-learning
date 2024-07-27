package controllers

import (
	"context"
	"fmt"
	"os"

	"github.com/StevenSermeus/go-learning/db_access"
	"github.com/StevenSermeus/go-learning/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type UserCreate struct {
	Username string `json:"username" binding:"required,min=2,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type Header struct {
	API_KEY string `header:"X-API-KEY" binding:"required,min=64,max=64"`
}


func CreateUser(c *gin.Context) {

	API_KEY := Header{}
	if c.BindHeader(&API_KEY) != nil {
		c.JSON(400, gin.H{"error": "API_KEY is required and should be 64 characters long"})
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

	app, err := queries.GetApplicationByAPIKey(ctx, API_KEY.API_KEY)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	user := UserCreate{}
	err = c.BindJSON(&user)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hash_pass, err := utils.GenerateFromPassword(user.Password, &utils.ArgonParams{ Memory: 2 * 1024, Iterations: 1, Parallelism: 2, SaltLength: 16, KeyLength: 32})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := queries.CreateUser(ctx, db_access.CreateUserParams{Username: user.Username, Pass: hash_pass, ApplicationID: app.ID})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//Convert bytes to string
	s := fmt.Sprintf("%x-%x-%x-%x-%x", createdUser.ID.Bytes[0:4], createdUser.ID.Bytes[4:6], createdUser.ID.Bytes[6:8], createdUser.ID.Bytes[8:10], createdUser.ID.Bytes[10:16])
	jwtToken, err := utils.CreateAccessToken(s,createdUser.Username, app.AccessTokenDuration)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	refresh_token , err := utils.CreateRefreshToken(s, createdUser.Username, app.RefreshTokenDuration)
 
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"last_login": createdUser.LastLogin, "username": createdUser.Username, "id": createdUser.ID, "access_token": jwtToken, "refresh_token": refresh_token})
}

type SignIn struct {
	Username string `json:"username" binding:"required,min=2,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type SignInHeader struct {
	API_KEY string `header:"X-API-KEY" binding:"required,min=64,max=64"`
}

func LoginHandler(c *gin.Context) {
	API_KEY := SignInHeader{}
	if c.BindHeader(&API_KEY) != nil {
		c.JSON(400, gin.H{"error": "API_KEY is required and should be 64 characters long"})
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

	app, err := queries.GetApplicationByAPIKey(ctx, API_KEY.API_KEY)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	signIn := SignIn{}
	err = c.BindJSON(&signIn)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := queries.GetUserByUsername(ctx, signIn.Username)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	match, err := utils.ComparePasswordAndHash(signIn.Password, user.Pass)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !match {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	if user.ApplicationID != app.ID {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}
	
	//Convert bytes to string
	s := fmt.Sprintf("%x-%x-%x-%x-%x", user.ID.Bytes[0:4], user.ID.Bytes[4:6], user.ID.Bytes[6:8], user.ID.Bytes[8:10], user.ID.Bytes[10:16])
	jwtToken, err := utils.CreateAccessToken(s,user.Username, app.AccessTokenDuration)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	refresh_token , err := utils.CreateRefreshToken(s, user.Username, app.RefreshTokenDuration)
 
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"last_login": user.LastLogin, "username": user.Username, "id": user.ID, "access_token": jwtToken, "refresh_token": refresh_token})
}