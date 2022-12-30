package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type dialog struct {
	Text  string `json:"text"`
	Sayer string `json:"sayer"`
}

func loadDialogs(data []byte) map[string][]dialog {
	var dialogMap map[string][]dialog

	if err := json.Unmarshal(data, &dialogMap); err != nil {
		panic("Can't unmarshal dialogs")
	}

	for key, chain := range dialogMap {
		dialogMap[key] = splitTexts(chain)
	}

	return dialogMap
}

func splitTexts(dialogs []dialog) []dialog {
	var modifiedDialogs []dialog
	for _, d := range dialogs {
		fmt.Println(d.Text)
		var splitted []string
		for i, char := range d.Text {
			if i%40 == 0 {
				splitted = append(splitted, "\n")
			}
			splitted = append(splitted, fmt.Sprintf("%c", char))
		}

		d.Text = strings.Join(splitted, "")
		fmt.Println(d.Text)

		modifiedDialogs = append(modifiedDialogs, d)
	}

	return modifiedDialogs
}
