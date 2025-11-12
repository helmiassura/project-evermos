package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"evermos-project/models"
	"evermos-project/utils"
)

const UserIDKey = "userID"

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// 1. Ambil token (dari header "token" atau "Authorization: Bearer <token>")
		tokenString := c.Get("token")

		if tokenString == "" {
			auth := c.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tokenString = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if tokenString == "" {
			return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Unauthorized", []string{"Token required"}, nil)
		}

		// 2. Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || token == nil || !token.Valid {
			return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Unauthorized", []string{"Invalid token"}, nil)
		}

		// 3. Ambil claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Unauthorized", []string{"Cannot parse token claims"}, nil)
		}

		// 4. Ambil user ID dari claim (harapkan ada key "id" yang bertipe number)
		idValue, ok := claims["id"]
		if !ok {
			return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Unauthorized", []string{"ID not found in token"}, nil)
		}

		userIDFloat, ok := idValue.(float64)
		if !ok {
			// sometimes id might be string in token; try string -> int
			if _, ok2 := idValue.(string); ok2 {
				// avoid importing strconv here; keep middleware simple
				return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Unauthorized", []string{"Invalid ID format in token (string) - expected number"}, nil)
			}
			return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Unauthorized", []string{"Invalid ID format"}, nil)
		}

		uid := uint(userIDFloat)
		c.Locals(UserIDKey, uid)    // main key used across project
		c.Locals("userID", uid)     // explicit
		c.Locals("user_id", uid)    // snake_case
		c.Locals("id", uid)         // older controllers might expect "id"

		return c.Next()
	}
}

func AdminMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawID := c.Locals(UserIDKey)
		if rawID == nil {
			return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Unauthorized", []string{"No user session"}, nil)
		}

		userID := rawID.(uint)

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Unauthorized", []string{"User not found"}, nil)
		}

		if !user.IsAdmin {
			return utils.RespondJSON(c, fiber.StatusForbidden, false, "Forbidden", []string{"Admin access required"}, nil)
		}

		return c.Next()
	}
}
