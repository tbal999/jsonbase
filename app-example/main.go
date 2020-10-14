package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tbal999/jsonbase"
)

func main() {
	d := jsonbase.Database{}
	quit := false
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("JBUI (q to quit, h for cmdlist)")
	for quit == false {
		fmt.Printf(">>> ")
		scanner.Scan()
		result := scanner.Text()
		switchy := strings.Split(result, " ")
		switch switchy[0] {
		case "q":
			quit = true
		case "h":
			fmt.Println(`
			command - info
			(for further info, just type in the command and further usage instruction will appear)

			count - lets you count number of instances of all unique items in column
			clear - clear all data in buffer
			describe - describe a table on one column
			datetodays - add new column to table converting a date column into a days from delta
			export - exports buffer as a csv file
			grab - grab a table of your choice to place in buffer
			import - imports a file of your choice into database. 
			join - left join one table against another
			knnclass - lets you do knn classification on dataset
			knnreg - lets you do knn regression on dataset
			load - load a database
			norm - normalize a dataset in a table of your choice
			order - order data in a table by one column either asc or desc
			plot - generate a line plot for table against one column
			print - prints out data in buffer
			row - lets you grab a specific row number of items that have already been ordered
			regex - grab data in table where column matches a regular expression
			regexrep - replace data in column of table where they match a regular expression
			savebuffer - saves current data in buffer as a new table
			sum - calculate sum of total in a column of integers
			shuffle - shuffle data in a table randomly
			save - save a database
			select - lets you select (and rename) columns in buffer
			split - lets you split up data into a training and testing set for machine learning
			unpivot - unpivots data in table on one column
			`)
		case "import":
			if len(switchy) == 5 {
				switch switchy[4] {
				case "true":
					d.ImportFile(switchy[1], switchy[2], switchy[3], true)
				case "false":
					d.ImportFile(switchy[1], switchy[2], switchy[3], false)
				}
			} else {
				fmt.Println("usage: import filepath tablename delimiter bool")
			}
		case "export":
			if len(switchy) == 2 {
				d.ExportAsCSV(switchy[1])
			} else {
				fmt.Println("usage: export file")
			}
		case "print":
			if len(switchy) == 3 {
				number, _ := strconv.Atoi(switchy[2])
				switch switchy[1] {
				case "true":
					d.Print(true, number)
				case "false":
					d.Print(false, number)
				}
			} else {
				fmt.Println("usage: print bool 0")
			}
		case "grab":
			if len(switchy) == 2 {
				d.GrabTable(switchy[1])
			} else {
				fmt.Println("usage: grab table")
			}
		case "clear":
			d.Clear()
		case "select":
			if len(switchy) != 1 {
				columns := []string{}
				tobeselected := switchy[1:]
				for index := range tobeselected {
					columns = append(columns, tobeselected[index])
				}
				d.Select(columns)
			} else {
				fmt.Println("usage: select column1[newname] column2[newname] etc")
			}
		case "save":
			if len(switchy) != 1 {
				d.SaveDBase(switchy[1])
			} else {
				fmt.Println("usage: save filename")
			}
		case "savebuffer":
			if len(switchy) != 1 {
				d.SaveBuffer(switchy[1], 0, false)
			} else {
				fmt.Println("usage: savebuffer newtablename")
			}
		case "load":
			if len(switchy) != 1 {
				d.LoadDBase(switchy[1])
			} else {
				fmt.Println("usage: load filename")
			}
		case "plot":
			if len(switchy) == 3 {
				d.Plot(switchy[1], switchy[2])
			} else {
				fmt.Println("usage: plot table column")
			}
		case "order":
			if len(switchy) == 4 {
				switch switchy[3] {
				case "asc":
					d.Order(switchy[1], switchy[2], true)
				case "desc":
					d.Order(switchy[1], switchy[2], false)
				}
			} else {
				fmt.Println("usage: order table column asc/desc")
			}
		case "row":
			if len(switchy) == 4 {
				row, _ := strconv.Atoi(switchy[3])
				d.Row(switchy[1], switchy[2], row)
			} else {
				fmt.Println("usage: row table column number")
			}
		case "norm":
			if len(switchy) == 2 {
				d.Normalize(switchy[1])
			} else {
				fmt.Println("usage: norm table")
			}
		case "join":
			if len(switchy) == 5 {
				d.Join(switchy[1], switchy[2], switchy[3], switchy[4])
			} else {
				fmt.Println("usage: join table1 column1 table2 column2")
			}
		case "split":
			if len(switchy) == 5 {
				ratio, _ := strconv.Atoi(switchy[4])
				d.Split(switchy[1], switchy[2], switchy[3], ratio)
			} else {
				fmt.Println("usage: split table newtable1 newtable2 number")
			}
		case "knnclass":
			if len(switchy) == 6 {
				var training bool
				if switchy[5] == "true" {
					training = true
				} else {
					training = false
				}
				knumber, _ := strconv.Atoi(switchy[4])
				d.KNNclass(switchy[1], switchy[2], switchy[3], knumber, training)
			} else {
				fmt.Println("usage: knnclass traintable testtable column knumber bool")
			}
		case "knnreg":
			if len(switchy) == 5 {
				knumber, _ := strconv.Atoi(switchy[4])
				d.KNNreg(switchy[1], switchy[2], switchy[3], knumber)
			} else {
				fmt.Println("usage: knnreg traintable testtable column knumber")
			}
		}
	}
}
