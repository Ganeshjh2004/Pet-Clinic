package middleware

import (
	"net/http"
	"petclinic/utils"
)

// RoleBasedAccess creates a middleware that allows access only if the user's role matches one of the allowedRoles
func RoleBasedAccess(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claimsVal := r.Context().Value("userClaims")
			claims, ok := claimsVal.(*utils.Claims)
			if !ok || claims == nil {
				http.Error(w, "Forbidden: no user claims", http.StatusForbidden)
				return
			}

			allowed := false // loop to check for the roles access
			for _, role := range allowedRoles {
				if claims.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Forbidden: insufficient role", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
