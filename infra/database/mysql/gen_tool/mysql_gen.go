package main

import (
	"os"

	"github.com/KAMIENDER/golang-scaffold/infra/config"
	"github.com/KAMIENDER/golang-scaffold/infra/database/mysql"
	"gorm.io/gen"
)

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := os.Setenv("config_path", workingDirectory+"/../../../../conf/dev.yaml"); err != nil {
		panic(err)
	}

	// 指定生成代码的具体(相对)目录，默认为：./query
	// 默认情况下需要使用WithContext之后才可以查询，但可以通过设置gen.WithoutContext避免这个操作
	g := gen.NewGenerator(gen.Config{
		// 最终package不能设置为model，在有数据库表同步的情况下会产生冲突，若一定要使用可以单独指定model package的新名字
		OutPath:      workingDirectory + "/../../../persistent/basic",
		ModelPkgPath: workingDirectory + "/../../../persistent/po", // 默认情况下会跟随OutPath参数，在同目录下生成model目录
		OutFile:      workingDirectory + "/../../../persistent/basic/basic_query.go",
		// WithUnitTest: true,
		/* Mode: gen.WithoutContext,*/
	})

	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	db, err := mysql.NewDatabase(conf)
	if err != nil {
		panic(err)
	}
	g.UseDB(db)

	dataMap := map[string]func(detailType string) (dataType string){
		"bigint":    func(detailType string) (dataType string) { return "int64" },
		"int":       func(detailType string) (dataType string) { return "int" },
		"bit":       func(detailType string) (dataType string) { return "bool" },
		"timestamp": func(detailType string) (dataType string) { return "*time.Time" },
	}

	g.WithDataTypeMap(dataMap)
	g.ApplyBasic(
		g.GenerateModelAs("user", "User"),
		g.GenerateModel("user", gen.FieldType("deleted_time", "gorm.DeletedAt")),
	)

	// gen.FieldType("deleted_at", "soft_delete.DeletedAt"),

	g.Execute()
	// gen.FieldGORMTag("deleted_at", "column:deleted_at;type:bigint(20) unsigned;not null;softDelete:milli;default:0"))
}
