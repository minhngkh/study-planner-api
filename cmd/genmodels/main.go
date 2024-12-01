package main

import (
	"strings"

	"gorm.io/gen"

	"study-planner-api/internal/db"
)

func main() {
	db := db.Get().DB

	g := gen.NewGenerator(gen.Config{
		ModelPkgPath: "internal/model",
	})

	g.WithTableNameStrategy(func(tableName string) string {
		if strings.HasPrefix(tableName, "_") || strings.HasPrefix(tableName, "sqlite_") {
			return ""
		}
		return tableName
	})

	g.UseDB(db)

	g.GenerateAllTable()

	g.Execute()
}
