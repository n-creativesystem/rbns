package restserver

const role = "admin"

// func roleCheck(compile *jmespath.JMESPath) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// string
// 		data, ok := c.Get(auth.AuthKey)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
// 			return
// 		}
// 		token, ok := data.(string)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
// 			return
// 		}
// 		var mp interface{}
// 		if err := json.Unmarshal([]byte(token), &mp); err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
// 			return
// 		}
// 		if result, err := compile.Search(mp); err != nil {
// 			logrus.Error(err)
// 		} else {
// 			if v, ok := result.(string); ok {
// 				if v == role {
// 					c.Next()
// 					return
// 				}
// 			}
// 		}
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
// 	}
// }
