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

//Golbal maps column header with it's respective column index
type hMap map[string]int
type xlsxPairs map[string]string

var pairs xlsxPairs

func findHeaders(sheets []*xlsx.Sheet) hMap {

	headers := make(hMap)

	for _, sheet := range sheets {
		if strings.EqualFold(sheet.Name, "Sheet1") {
			for _, row := range sheet.Rows[0:1] {
				for idx, cell := range row.Cells {
					headers[strings.ToLower(cell.Value)] = idx
				}
			}
		}
	}

	log.Println(headers)

	return headers
}

func findSheets(header map[string]int, sheets []*xlsx.Sheet) xlsxPairs {

	xP := make(xlsxPairs)

	keyHeaders := []string{"key", "value"}
	log.Println("Header with column Index:", header)

	for _, sheet := range sheets {
		for _, row := range sheet.Rows[1:] {
			if colIdx, ok := header[keyHeaders[0]]; ok && len(row.Cells) > colIdx {
				key := strings.TrimSpace(row.Cells[colIdx].Value)
				if colIdx, ok := header[keyHeaders[1]]; ok && len(row.Cells) > colIdx {
					if key != "" {
						value := strings.TrimSpace(row.Cells[colIdx].Value)
						xP[key] = value
					}
				}
			}
		}
	}

	return xP
}

func init() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := dir + "/controllers/file.xlsx"
	fmt.Println("File  path:", path)

	xLSXFile, err := xlsx.OpenFile(path)
	if err != nil {
		panic(err)
		return
	}

	pairs = findSheets(findHeaders(xLSXFile.Sheets), xLSXFile.Sheets)
}

//ProcessXLSX  handler takes the key input as rest parameters
//Response is the value from the target XLSX sheet
func ProcessXLSX(w http.ResponseWriter, r *http.Request) {

	params := r.Context().Value("params").(httprouter.Params)
	key := params.ByName("key")

	log.Println(key)

	res := make(map[string]interface{})
	if key == "" {
		res["response"] = "Invalid entry"
		renderJSON(w, http.StatusOK, res)
		return
	}

	type response struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	resObj := response{}
	resObj.Key = key
	resObj.Value = pairs[key]

	if resObj.Value == "" {
		resObj.Value = "N/A"
		res["response"] = "Invalid entry"
		renderJSON(w, http.StatusOK, res)
		return
	}

	res["response"] = resObj
	renderJSON(w, http.StatusOK, res)
}
