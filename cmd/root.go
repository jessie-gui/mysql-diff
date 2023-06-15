package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mysql-diff/core"
	"mysql-diff/helper"
)

func init() {
	rootCmd.Flags().StringVarP(&source, "source", "s", "", "源服务器(格式: <user>:<password>@<host>:<port>)")
	rootCmd.Flags().StringVarP(&target, "target", "t", "", "目标服务器(格式: <user>:<password>@<host>:<port>)")
	rootCmd.Flags().StringVarP(&db, "db", "d", "", "指定数据库(格式: <source_db>:<target_db>)")
}

var (
	source, target, db string

	rootCmd = &cobra.Command{
		Use:     "mysqldiff",
		Short:   "mysql差异比较工具",
		Version: "v1.0.0",
		Run: func(cmd *cobra.Command, args []string) {
			config := helper.NewConfig()

			// 若没指定服务器，则默认用配置文件中的服务器配置
			if source == "" {
				source = fmt.Sprintf("%s:%s@%s", config.Source.User, config.Source.Pwd, config.Source.Address)
			}

			if target == "" {
				target = fmt.Sprintf("%s:%s@%s", config.Target.User, config.Target.Pwd, config.Target.Address)
			}

			if db == "" {
				db = fmt.Sprintf("%s:%s", config.Source.DbBase, config.Target.DbBase)
			}

			// 校验参数格式
			mysqlMap := core.NewParams(core.Source(source), core.Target(target), core.Db(db)).CheckParams()

			// 获取数据库连接
			sourceConn := mysqlMap["sourceMysql"].Conn()
			targetConn := mysqlMap["targetMysql"].Conn()

			// 开始比较
			core.NewCompare(sourceConn, targetConn).Compare(mysqlMap["sourceMysql"], mysqlMap["targetMysql"])
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}
