package docs

import (
	httpSwagger "github.com/swaggo/http-swagger"
)

/* Инициализация Swagger */
func init() {
	httpSwagger.PersistAuthorization(true)
}
