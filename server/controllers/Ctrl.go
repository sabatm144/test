package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/tealeg/xlsx"
)

type headerMap map[string]int

func ReadXLSX(w http.ResponseWriter, r *http.Request) {

	params := r.Context().Value("params").(httprouter.Params)
	key := params.ByName("key")

	log.Println(key)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := dir + "/controllers/file.xlsx"
	fmt.Println("Directory 1:", path)
	fmt.Println("Directory :", path)

	xLSXFile, err := xlsx.OpenFile(path)
	if err != nil {
		panic(err)
		return
	}

	data := headerMap{}
	headers := []string{"key", "value"}
	xlsx_pairs := map[string]string{}
	data["key"] = 0
	data["value"] = 1

	for _, sheet := range xLSXFile.Sheets {
		for _, row := range sheet.Rows[1:] {

			if colIdx, ok := data[headers[0]]; ok && len(row.Cells) > colIdx {
				key := strings.TrimSpace(row.Cells[colIdx].Value)

				if colIdx, ok := data[headers[1]]; ok && len(row.Cells) > colIdx {
					value := strings.TrimSpace(row.Cells[colIdx].Value)
					if key != "" {
						xlsx_pairs[key] = value
					}
				}
			}
		}
	}

	res := make(map[string]interface{})

	type response struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	
	resObj := response{}
	resObj.Key = key
	resObj.Value = xlsx_pairs[key]

	if resObj.Value == "" {
		res["response"] = "Invalid key"
		renderJSON(w, http.StatusOK, res)
		return
	}
	res["response"] = resObj
	renderJSON(w, http.StatusOK, res)
}
