package jwt

import (
	"time"

	"github.com/shevchenko-a-v/auth-service/internal/domain/models"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {

}
