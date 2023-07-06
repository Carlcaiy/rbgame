package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func TestExcel(t *testing.T) {
	f, err := excelize.OpenFile("slot.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows := f.GetRows("slot")
	sum, _ := strconv.Atoi(rows[0][4])

	count := 0
	value := 0
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if row[2] != "95" && row[2] != "1039" {
			fmt.Println(row)
			continue
		}

		change, err := strconv.Atoi(row[3])
		if err != nil {
			fmt.Println(err)
			break
		}
		target, err := strconv.Atoi(row[4])
		if err != nil {
			fmt.Println(err)
			break
		}
		sum += change
		if row[2] == "95" {
			count++
			if count == 1 {
				value = target
			}
			continue
		} else if count > 1 {
			target = value
		}
		count = 0
		if sum != target {
			fmt.Println(i+1, sum, change, target)
			break
		}
	}
	fmt.Println(count)
}
