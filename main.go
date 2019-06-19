package main

import (
	"encoding/json"
	"flag"
	"github.com/savaki/jq"
	"github.com/jmoiron/jsonq"
	"github.com/jedib0t/go-pretty/table"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	rowLocation = flag.String("row", ".hits.hits", "Row location https://github.com/savaki/jq query path")
	columnsFlag = flag.String("columns", "_id:id,_source.fieldName1:fieldName1,_source.fieldName2:fieldName2,_source.fieldName3:fieldName3", "Field order and selection")
)

func main() {
	log.SetFlags(log.Flags()|log.Lshortfile)
	flag.Parse()
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Panic(err)
	}
	rowQuery, err := jq.Parse(*rowLocation)
	if err != nil {
		log.Panic(err)
	}
	rowsQ, err := rowQuery.Apply(b)
	if err != nil {
		log.Panic(err)
	}
	rows := []map[string]interface{}{}
	if err := json.Unmarshal(rowsQ, &rows); err != nil {
		log.Printf("%#v", string(rowsQ))
		log.Panic(err)
	}
	columns := strings.Split(*columnsFlag, ",")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	headerRow := table.Row{}
	for _, column := range columns {
		columnName := column
		columnSplit := strings.Split(column, ":")
		if len(columnSplit) > 1 {
			columnName = columnSplit[1]
		}
		headerRow = append(headerRow, columnName)
	}
	t.AppendHeader(headerRow)

	for _, rowMap := range rows {
		row := jsonq.NewQuery(rowMap)
		tableRow := table.Row{}
		for _, column := range columns {
			s, err := row.String(strings.Split(strings.Split(column, ":")[0], ".")...)
			if err != nil {
				log.Printf("Error: %v", err)
			}
			tableRow = append(tableRow, s)
		}
		t.AppendRow(tableRow)
	}
	t.RenderMarkdown()

}
