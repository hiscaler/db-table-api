package main

import (
	"log"
	"net/http"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/slash"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"strconv"
)

func main() {
	router := routing.New()
	router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	api := router.Group("/api")
	api.Use(
		content.TypeNegotiator(content.JSON),
	)

	db, _ := dbx.Open("mysql", "root:root@/touch_admin_s")

	// GET /api/
	api.Get("/", func(c *routing.Context) error {
		return c.Write("OK")
	})

	// GET /api/TABLE_NAME?page=1&pageSize=100
	api.Get(`/<table:\w+>`, func(c *routing.Context) error {
		page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
		pageSize, _ := strconv.ParseInt(c.Query("pageSize", "100"), 10, 64)
		table := c.Param("table")
		table = strings.Trim(table, " ")
		if len(table) == 0 {
			panic("Table name is can't empty.")
		}
		res :=make([]interface{}, 0)
		row := dbx.NullStringMap{}
		rows, _ := db.Select().From(table).Offset((page - 1) * pageSize).Limit(pageSize).Rows()
		for rows.Next() {
			rows.ScanMap(row)
			t := make(map[string]interface{})
			for name, v := range row {
				t[name] = v.String
			}
			res = append(res, t)
		}
		return c.Write(res)
	})

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
