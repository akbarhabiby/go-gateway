# go-gateway

A Gateway Server
 
**Basic Samples:**
- Gateway Server PORT: 5000

`Sample 1`

  **Code:** 
  ```go
  proxyConfigs := []proxy.ProxyMiddlewareConfig{
		{
			Name:    "buntil",
			Prefix:  "/buntil/callback",
			URL:     "http://localhost:3000",
			Rewrite: map[string]string{"/buntil/callback": "/callback"},
		},
	}
  ```
  **Description:** 
  - This will rewrite from middleware server "http://localhost:5000/buntil/callback" to "http://localhost:3000/callback"

---

`Sample 2`

  **Code:** 
  ```go
  proxyConfigs := []proxy.ProxyMiddlewareConfig{
		{
			Name:    "buntil",
			Prefix:  "/test",
			URL:     "http://localhost:3001",
			Timeout: time.Second * 5,
      Rewrite: map[string]string{"/test": "/"},
		},
	}
  ```
  **Description:** 
  - This will rewrite from middleware server "http://localhost:5000/test" to "http://localhost:3000/"

---


`Sample 3`

  **Code:** 
  ```go
  proxyConfigs := []proxy.ProxyMiddlewareConfig{
		{
			Prefix:  "/test",
			URL:     "http://localhost:3001",
		},
	}
  ```
  **Description:** 
  - This will rewrite from middleware server "http://localhost:5000/test" to "http://localhost:3001/test"
