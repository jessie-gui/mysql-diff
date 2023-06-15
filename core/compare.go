package core

import (
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"sort"
	"sync"
)

var (
	diffSqlKeys []string
	diffSqlMap  = make(map[string]string)
	wg          sync.WaitGroup
	ch          = make(chan bool, 16)
	lock        sync.Mutex
	comment     bool
	foreign     bool
	tidb        bool
)

// Compare /**
type Compare struct {
	sourceDb *gorm.DB
	targetDb *gorm.DB
}

func NewCompare(sourceDb, targetDb *gorm.DB) *Compare {
	return &Compare{
		sourceDb: sourceDb,
		targetDb: targetDb,
	}
}

func (c *Compare) Compare(sourceMysql, targetMysql *Mysql) {
	var (
		sourceSchema Schema
		targetSchema Schema
	)

	sourceSchemaResult := c.sourceDb.Table("SCHEMATA").Limit(1).Find(
		&sourceSchema,
		"`SCHEMA_NAME` = ?", sourceMysql.GetDb(),
	)

	targetSchemaResult := c.targetDb.Table("SCHEMATA").Limit(1).Find(
		&targetSchema,
		"`SCHEMA_NAME` = ?", targetMysql.GetDb(),
	)

	if sourceSchemaResult.RowsAffected <= 0 {
		cobra.CheckErr(fmt.Errorf("源数据库 `%s` 不存在。", sourceMysql.GetDb()))
	}

	if targetSchemaResult.RowsAffected <= 0 {
		cobra.CheckErr(fmt.Errorf("目标数据库 `%s` 不存在。", targetMysql.GetDb()))
	}

	var (
		sourceTableData []Table
		targetTableData []Table
	)

	c.sourceDb.Table("TABLES").Order("`TABLE_NAME` ASC").Find(
		&sourceTableData,
		"`TABLE_SCHEMA` = ?", sourceMysql.GetDb(),
	)
	c.targetDb.Table("TABLES").Order("`TABLE_NAME` ASC").Find(
		&targetTableData,
		"`TABLE_SCHEMA` = ?", targetMysql.GetDb(),
	)

	sourceTableMap := make(map[string]Table)
	targetTableMap := make(map[string]Table)

	for _, table := range sourceTableData {
		sourceTableMap[table.TableName] = table
	}

	for _, table := range targetTableData {
		targetTableMap[table.TableName] = table
	}

	// DROP TABLE Or DROP VIEW...
	drop(sourceTableMap, targetTableData)

	defer close(ch)

	for _, sourceTable := range sourceTableData {
		wg.Add(1)

		go diff(sourceMysql, targetMysql, c.sourceDb, c.targetDb, sourceSchema, sourceTable, targetTableMap)
	}

	wg.Wait()

	// Print Sql...
	if len(diffSqlKeys) > 0 && len(diffSqlMap) > 0 {
		fmt.Println(fmt.Sprintf("SET NAMES %s;\n", sourceSchema.DefaultCharacterSetName))
		fmt.Println("SET FOREIGN_KEY_CHECKS=0;")
		fmt.Println()

		sort.Strings(diffSqlKeys)

		for k, diffSqlKey := range diffSqlKeys {
			if diffSql, ok := diffSqlMap[diffSqlKey]; ok {
				if k < len(diffSqlKeys)-1 {
					fmt.Println(diffSql)
					fmt.Println()
				} else {
					fmt.Println(diffSql)
				}
			}
		}

		fmt.Println()
		fmt.Println("SET FOREIGN_KEY_CHECKS=1;")
	} else {
		fmt.Println("数据库无差异！")
	}
}
