package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

type Client struct {
	Requests int
	ResetAt  time.Time
}

var (
	clients = make(map[string]*Client)
	mu      sync.Mutex
	limit   = 100
	window  = 1 * time.Minute
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getClientKey(c)
		now := time.Now()

		mu.Lock()
		client, exists := clients[key]

		if !exists || now.After(client.ResetAt) {
			clients[key] = &Client{
				Requests: 1,
				ResetAt:  now.Add(window),
			}
			mu.Unlock()
			c.Next()
			return
		}

		client.Requests++

		if client.Requests > limit {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		mu.Unlock()
		c.Next()
	}
}

func getClientKey(c *gin.Context) string {
	authCookie, err := c.Cookie("access_token")

	if err != nil || authCookie == "" {
		return c.ClientIP()
	}

	token, _, err := new(jwt.Parser).ParseUnverified(authCookie, jwt.MapClaims{})
	if err != nil {
		return c.ClientIP()
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userId, exists := claims["user_id"]; exists {
			if idStr, ok := userId.(string); ok {
				return idStr
			}
		}
	}

	return c.ClientIP()
}

func StartRateLimitCleanup(lc fx.Lifecycle) {
	ticker := time.NewTicker(5 * time.Minute)
	stop := make(chan struct{})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				for {
					select {
					case <-ticker.C:
						cleanupClients()
					case <-stop:
						ticker.Stop()
						return
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			close(stop)
			return nil
		},
	})
}

func cleanupClients() {
	now := time.Now()

	mu.Lock()
	defer mu.Unlock()

	for key, client := range clients {
		if now.After(client.ResetAt.Add(window)) {
			delete(clients, key)
		}
	}
}
