package main

import (
	"log"

	"github.com/jessie-gui/mysql-diff/cmd"
)

/**
 *
 *
 * @author        Gavin Gui <guijiaxian@gmail.com>
 * @version       1.0.0
 * @copyright (c) 2022, Gavin Gui
 */
func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
