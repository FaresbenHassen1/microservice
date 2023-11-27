package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// Create an HTTP gin middleware to validate the bearer token
func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		provider, err := oidc.NewProvider(context.Background(), "http://127.0.0.1:5556/dex")
		if err != nil {
			log.Fatal(err)
		}
		verifier := provider.Verifier(&oidc.Config{ClientID: "example-app"})

		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		tokenString := parts[1]

		token, err := verifier.Verify(context.Background(), tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
		var UserInfo struct {
			Subject  string `json:"sub"`
			Email    string `json:"email"`
			Verified bool   `json:"email_verified"`
			Name     string `json:"name"`
		}
		if err := token.Claims(&UserInfo); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
		fmt.Println(UserInfo)
		fmt.Println("Subject :", UserInfo.Subject)
		fmt.Println("Email :", UserInfo.Email)
		fmt.Println("Email's verification :", UserInfo.Verified)
		fmt.Println("Name :", UserInfo.Name)
		// ctx.Set("sub", claims)//
		ctx.Next()
	}
}
