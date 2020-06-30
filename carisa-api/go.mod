module github.com/carisa/api

go 1.14

replace github.com/carisa/pkg => ../carisa-pkg

require (
	github.com/carisa/pkg v0.0.0-00010101000000-000000000000
	github.com/labstack/echo/v4 v4.1.16
	github.com/rs/xid v1.2.1
	github.com/stretchr/testify v1.4.0
	gopkg.in/yaml.v2 v2.2.2
)
