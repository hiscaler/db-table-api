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
	"dta/response"
)

var cfg = struct {
	Debug           bool
	ListenPort      string
	Driver          string
	DSN             string
	database        string
	tables          []string
	TablePrefix     string
	FieldNameFormat string
	toCamel         bool
}{
	Debug:           true,
	Driver:          "mysql",
	FieldNameFormat: "original",
	toCamel:         false,
}

var db *dbx.DB

// 转换字符串为驼峰格式
func toCamel(s, sep string) string {
	if len(sep) == 0 {
		return s
	}

	names := strings.Split(s, sep)
	vv := make([]string, len(names))
	for _, v := range names {
		vv = append(vv, strings.ToUpper(v[:1])+strings.ToLower(v[1:]))
	}
	return strings.Join(vv, "")
}

func parseTable(table string) string {
	table = strings.Trim(table, " ")
	if len(table) == 0 {
		panic("Table name is can't empty.")
	}

	if len(cfg.TablePrefix) > 0 && !strings.HasPrefix(table, cfg.TablePrefix) {
		table = cfg.TablePrefix + table
	}

	if len(cfg.tables) != 0 {
		tableValid := false
		for _, v := range cfg.tables {
			if v == table {
				tableValid = true
				break
			}
		}

		if !tableValid {
			panic(fmt.Sprintf("Table `%s` not exists in database.", table))
		}
	}

	return table
}

func init() {
	// Read and get app config
	jsonFile, err := ioutil.ReadFile("src/dta/config/conf.json")
	if err != nil {
		panic("Read config file failed.")
	}

	err = json.Unmarshal(jsonFile, &cfg)
	if err != nil {
		panic("Config file is invalid. Must be a valid json format.");
	}

	if strings.ToLower(strings.Trim(cfg.FieldNameFormat, " ")) == "camel" {
		cfg.toCamel = true
	}

	t := strings.Split(cfg.DSN, "/")
	if len(t) < 2 {
		panic("DSN error.")
	} else {
		db, _ = dbx.Open(cfg.Driver, cfg.DSN)
		cfg.database = t[1]
		if strings.ToLower(cfg.Driver) == "mysql" && len(cfg.tables) == 0 {
			// Only support MySQL
			db.NewQuery(fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s'", cfg.database)).Column(&cfg.tables)
		}
	}
}

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

	// GET /api/
	api.Get("/ping", func(c *routing.Context) error {
		return c.Write("OK")
	})

	// 获取数据列表
	// GET /api/TABLE_NAME?page=1&pageSize=100
	api.Get(`/<table:\w+>`, func(c *routing.Context) error {
		page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
		pageSize, _ := strconv.ParseInt(c.Query("pageSize", "100"), 10, 64)
		table := parseTable(c.Param("table"))

		row := dbx.NullStringMap{}
		q := db.Select().From(table)
		var totalCount int64
		q.Select("COUNT(*)").Row(&totalCount)
		totalPages := (totalCount + pageSize - 1) / pageSize
		rows, err := q.Select("*").Offset((page - 1) * pageSize).Limit(pageSize).Rows()
		if err == nil {
			d := &response.SuccessListData{}
			d.Items = make([]interface{}, 0)

			for rows.Next() {
				rows.ScanMap(row)
				t := make(map[string]interface{})
				for name, v := range row {
					if cfg.toCamel {
						name = toCamel(name, "_")
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
			resp := &response.SuccessListResponse{
				Success: true,
				Data:    *d,
			}
			return c.Write(resp)
		} else {
			error := &response.Error{
				Message: fmt.Errorf("%v", err).Error(),
			}
			resp := &response.FailResponse{
				Success: false,
				Error:   *error,
			}
			return c.Write(resp)
		}
	})

	// 获取指定的数据
	// GET /api/TABLE_NAME/ID
	api.Get(`/<table:\w+>/<id:\d+>`, func(c *routing.Context) error {
		table := parseTable(c.Param("table"))
		id := c.Param("id")
		row := dbx.NullStringMap{}
		err := db.Select().From(table).Where(dbx.HashExp{"id": id}).One(row)
		if err == nil {
			data := make(map[string]interface{})
			for name, v := range row {
				if cfg.toCamel {
					name = toCamel(name, "_")
				}
				data[name] = v.String
			}
			resp := &response.SuccessOneResponse{
				Success: true,
				Data:    data,
			}
			return c.Write(resp)
		} else {
			error := &response.Error{
				Message: fmt.Errorf("%v", err).Error(),
			}
			resp := &response.FailResponse{
				Success: false,
				Error:   *error,
			}
			return c.Write(resp)
		}

	})

	http.Handle("/", router)
	addr := cfg.ListenPort
	if len(addr) == 0 {
		addr = "8080"
	}
	http.ListenAndServe(":"+addr, nil)
}
