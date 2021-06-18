package middleware

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

// Authorize determines if current subject has been authorized to take an action on an object.
func Authorize(obj string, act string, adapter *gormadapter.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		val, existed := c.Get("userID")
		log.Print("current User", val)
		if !existed {
			c.AbortWithStatusJSON(401, gin.H{"msg": "user hasn't logged in yet"})
			return
		}

		// Casbin enforces policy
		ok, err := enforce(val.(float64), obj, act, adapter)
		if err != nil {
			log.Print("err", err)
			c.AbortWithStatusJSON(500, gin.H{"msg": "error occurred when authorizing user"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"msg": "forbidden"})
			log.Print("second err")
			return
		}
		c.Next()
	}
}

func enforce(sub float64, obj string, act string, adapter *gormadapter.Adapter) (bool, error) {
	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return false, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// Load policies from DB dynamically
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}

	// Verify
	ok, err := enforcer.Enforce(fmt.Sprint(sub), obj, act)
	return ok, err
}
