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
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func main() {
	// Read and get app config
	jsonFile, err := ioutil.ReadFile("src/dta/config/conf.json")
	if err != nil {
		panic("Read config file failed.")
	}
	cfg := struct {
		Debug           bool
		Driver          string
		DSN             string
		TablePrefix     string
		FieldNameFormat string
	}{
		Debug:           true,
		Driver:          "mysql",
		FieldNameFormat: "original",
	}
	err = json.Unmarshal(jsonFile, &cfg)
	if err != nil {
		panic("Config file is invalid. Must be a valid json format.");
	}

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

	db, _ := dbx.Open(cfg.Driver, cfg.DSN)

	// GET /api/
	api.Get("/", func(c *routing.Context) error {
		return c.Write("OK")
	})

	// GET /api/TABLE_NAME?page=1&pageSize=100
	api.Get(`/<table:\w+>`, func(c *routing.Context) error {
		type Data struct {
			Items []interface{}    `json:"items"`
			Meta  map[string]int64 `json:"_meta"`
		}

		type Error struct {
			Message string `json:"message"`
		}

		type SuccessResponse struct {
			Success bool `json:"success"`
			Data    Data `json:"data"`
		}

		type FailResponse struct {
			Success bool  `json:"success"`
			Error   Error `json:"error"`
		}

		page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
		pageSize, _ := strconv.ParseInt(c.Query("pageSize", "100"), 10, 64)
		table := c.Param("table")
		table = strings.Trim(table, " ")
		if len(table) == 0 {
			panic("Table name is can't empty.")
		}

		if len(cfg.TablePrefix) > 0 && !strings.HasPrefix(table, cfg.TablePrefix) {
			table = cfg.TablePrefix + table
		}

		fieldNameCamelFormat := strings.ToLower(cfg.FieldNameFormat) == "camel"
		row := dbx.NullStringMap{}
		q := db.Select().From(table)
		var totalCount int64
		q.Select("COUNT(*)").Row(&totalCount)
		totalPages := (totalCount + pageSize - 1) / pageSize
		rows, err := q.Select("*").Offset((page - 1) * pageSize).Limit(pageSize).Rows()
		if err == nil {
			d := &Data{}
			d.Items = make([]interface{}, 0)

			for rows.Next() {
				rows.ScanMap(row)
				t := make(map[string]interface{})
				for name, v := range row {
					if fieldNameCamelFormat {
						names := strings.Split(name, "_")
						vv := make([]string, len(names))
						for _, v := range names {
							vv = append(vv, strings.ToUpper(v[:1])+strings.ToLower(v[1:]))
						}
						name = strings.Join(vv, "")
					}
					t[name] = v.String
				}
				d.Items = append(d.Items, t)
			}

			d.Meta = map[string]int64{
				"totalCount":  totalCount,
				"pageCount":   totalPages,
				"currentPage": page,
				"perPage":     pageSize,
			}
			resp := &SuccessResponse{
				Success: true,
				Data:    *d,
			}
			return c.Write(resp)

		} else {
			error := &Error{
				Message: fmt.Errorf("%v", err).Error(),
			}
			resp := &FailResponse{
				Success: false,
				Error:   *error,
			}
			return c.Write(resp)
		}
	})

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
