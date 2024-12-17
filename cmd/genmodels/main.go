package main

import (
	"strings"

	"gorm.io/gen"

	db "study-planner-api/internal/database"
)

func main() {
	db := db.Instance().DB

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
	// g.GenerateModel("task",
	// 	gen.FieldType("start_time", "time.Time"),
	// 	gen.FieldType("end_time", "time.Time"),
	// 	gen.FieldType("created_at", "time.Time"),
	// 	gen.FieldType("estimated_time", "int32"),
	// )

	g.Execute()
}
