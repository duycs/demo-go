package gsheet

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

func Run() {
	data, err := ioutil.ReadFile("client_secret.json")
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet("1eR-eOgG09RG8k8JDZWoPDoU0ovnsBTMvJ7DC5d3kcwE")
	checkError(err)
	sheet, err := spreadsheet.SheetByIndex(0)
	checkError(err)
	for _, row := range sheet.Rows {
		for _, cell := range row {
			fmt.Println(cell.Value)
		}
	}

	// Update cell content
	sheet.Update(0, 0, "update done")

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
