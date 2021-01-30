package jsonbase

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	models "github.com/tbal999/jsonbase/models"
	table "github.com/tbal999/jsonbase/table"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

//Database is a struct that stores all tables and one neural network (if you want to save and use ANN later on)
type Database struct {
	Table []table.Table
	Net models.Network
}

//Temptable is the internal buffer for storing any queries - printable via 'Output' function.
//Exported just in case you want to access it directly.
var Temptable [][]string

/////////////////////      HELPER FUNCTIONS  /////////////////////////////

func grabsubstring(x string) (string, string) {
	if yes, _ := regexp.MatchString(`\[([^\[\]]*)\]`, x); yes == true {
		re := regexp.MustCompile(`\[([^\[\]]*)\]`)
		submatchall := re.FindAllString(x, -1)
		for _, element := range submatchall {
			element = strings.Trim(element, "[")
			element = strings.Trim(element, "]")
			return element, strings.Split(x, "[")[0]
		}
	}
	return "", strings.Split(x, "[")[0]
}

/////////////////////      HELPER FUNCTIONS  /////////////////////////////

//checks if table exists in database
func (d Database) verifytable(name string) (bool, int) {
	for index := range d.Table {
		if d.Table[index].Name == name {
			return true, index
		}
	}
	fmt.Println("Table " + name + " does not exist.")
	return false, 0
}

//ImportFile lets you import delimited flat files. filename is the name of file, delimiter is the delimiter that the file is delimited by.
//Set 'header' to true if there's a header for the file, otherwise set to false.
//Header is only for rune delimited files i.e a comma, it doesn't matter what you set it to if the file is delimited by '\n'
func (d *Database) ImportFile(filename, tablename, delimiter string, header bool) {
	D := *d
	Table := table.Table{}
	if filename == "" || tablename == "" {
		fmt.Println("Need a file / tablename.")
		return
	}
	if delimiter != "\n" && delimiter != "\r\n" {
		f, _ := os.Open(filename)
		r := csv.NewReader(f)
		r.Comma = rune(delimiter[0])
		r.LazyQuotes = true
		r.FieldsPerRecord = -1
		count := 0
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}
			for index := range record {
				record[index] = strings.Replace(record[index], "\"", "", -1)
				record[index] = strings.Replace(record[index], ",", "", -1)
			}
			if header == true {
				if count == 0 {
					Table.Columns = record
				} else {
					Table.Rows = append(Table.Rows, record)
				}
			} else {
				Table.Rows = append(Table.Rows, record)
			}
			count++
		}
		if header == false {
			head := ""
			for x := 0; x < len(Table.Rows[0]); x++ {
				head += "COLUMN" + strconv.Itoa(x)
				if x != len(Table.Rows[0])-1 {
					head += ","
				}
			}
			out := strings.Split(head, ",")
			for index := range out {
				Table.Columns = append(Table.Columns, out[index])
			}
		}
	} else {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		s := strings.Split(string(content), delimiter)
		Columns := "COLUMN1"
		Table.Columns = append(Table.Columns, Columns)
		for index := range s {
			rows := []string{s[index]}
			Table.Rows = append(Table.Rows, rows)
		}
	}
	Table.Name = tablename
	D.Table = append(D.Table, Table)
	*d = D
}

//ImportString lets you import strings delimited by EOL. filename is the name of file, delimiter is the delimiter that the string array is delimited by.
//header is true or false depending on whether there's columns in the data, if not, auto generated columns will be added.
func (d *Database) ImportString(input string, tablename, delimiter string, header bool) {
	D := *d
	Table := table.Table{}
	stage1 := strings.Split(input, "\n")
	Table.Name = tablename
	for index := range stage1 {
		if header == true {
			if index == 0 {
				Table.Columns = append(Table.Columns, stage1[index])
			} else {
				Table.Rows = append(Table.Rows, strings.Split(stage1[index], delimiter))
			}
		} else {
			Table.Rows = append(Table.Rows, strings.Split(stage1[index], delimiter))
		}
	}
	if header == false {
		head := ""
		for x := 0; x < len(Table.Rows[0]); x++ {
			head += "COLUMN" + strconv.Itoa(x)
			if x != len(Table.Rows[0])-1 {
				head += ","
			}
		}
		mid := strings.Split(head, ",")
		for index := range mid {
			Table.Columns = append(Table.Columns, mid[index])
		}
	}
	D.Table = append(D.Table, Table)
	*d = D
}

//Import1DString lets you import 1D string arrays. filename is the name of file, delimiter is the delimiter that the string array is delimited by.
//header is true or false depending on whether there's columns in the data, if not, auto generated columns will be added.
func (d *Database) Import1DString(input []string, tablename, delimiter string, header bool) {
	D := *d
	Table := table.Table{}
	Table.Name = tablename
	for index := range input {
		if header == true {
			if index == 0 {
				Table.Columns = append(Table.Columns, input[index])
			} else {
				Table.Rows = append(Table.Rows, strings.Split(input[index], delimiter))
			}
		} else {
			Table.Rows = append(Table.Rows, strings.Split(input[index], delimiter))
		}
	}
	if header == false {
		head := ""
		for x := 0; x < len(Table.Rows[0]); x++ {
			head += "COLUMN" + strconv.Itoa(x)
			if x != len(Table.Rows[0])-1 {
				head += ","
			}
		}
		mid := strings.Split(head, ",")
		for index := range mid {
			Table.Columns = append(Table.Columns, mid[index])
		}
	}
	D.Table = append(D.Table, Table)
	*d = D
}

//Import2DString lets you import 2D string arrays. filename is the name of file, delimiter is the delimiter that the string array is delimited by.
//header is true or false depending on whether there's columns in the data, if not, auto generated columns will be added.
func (d *Database) Import2DString(input [][]string, tablename, delimiter string, header bool) {
	D := *d
	Table := table.Table{}
	Table.Name = tablename
	for index := range input {
		tobeadded := strings.Join(input[index], delimiter)
		if header == true {
			if index == 0 {
				Table.Columns = append(Table.Columns, tobeadded)
			} else {
				Table.Rows = append(Table.Rows, strings.Split(tobeadded, delimiter))
			}
		} else {
			Table.Rows = append(Table.Rows, strings.Split(tobeadded, delimiter))
		}
	}
	if header == false {
		head := ""
		for x := 0; x < len(Table.Rows[0]); x++ {
			head += "COLUMN" + strconv.Itoa(x)
			if x != len(Table.Rows[0])-1 {
				head += ","
			}
		}
		mid := strings.Split(head, ",")
		for index := range mid {
			Table.Columns = append(Table.Columns, mid[index])
		}
	}
	D.Table = append(D.Table, Table)
	*d = D
}

//Print prints out the Temptable, bool lets you determine if table is cleared after print.
//Howmany is how many rows you want printed (0 for all rows)
func (d Database) Print(clear bool, howmany int) {
	max := 0
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	if len(Temptable) == 0 {
		fmt.Println("No results.")
	}
	if howmany == 0 {
		for index := range Temptable {
			out := strings.Join(Temptable[index], "\t")
			fmt.Fprintln(w, out)
		}
	} else {
		for index := range Temptable {
			out := strings.Join(Temptable[index], "\t")
			fmt.Fprintln(w, out)
			max++
			if max > howmany {
				break
			}
		}
	}
	w.Flush()
	if clear == true {
		Temptable = [][]string{}
	}
}

//Transpose flips the Temptable so columns are rows and rows are columns.
func (d Database) Transpose() {
    xl := len(Temptable[0])
    yl := len(Temptable)
    result := make([][]string, xl)
    for i := range result {
        result[i] = make([]string, yl)
    }
    for i := 0; i < xl; i++ {
        for j := 0; j < yl; j++ {
            result[i][j] = Temptable[j][i]
        }
    }
    Temptable = result
}

//SaveBuffer lets you save the current Temptable as a jsonbase table - name is the name of the table, howmany is how many rows you want to save
//clear is whether you want to clear the buffer after you've saved it.
func (d *Database) SaveBuffer(name string, howmany int, clear bool) {
	D := *d
	max := 0
	Table := table.Table{}
	Table.Name = name
	if len(Temptable) == 0 {
		fmt.Println("No results.")
		return
	}
	if howmany == 0 {
		for index := range Temptable {
			if index == 0 {
				Table.Columns = Temptable[index]
			} else {
				Table.Rows = append(Table.Rows, Temptable[index])
			}
		}
	} else {
		for index := range Temptable {
			if index == 0 {
				Table.Columns = Temptable[index]
			} else {
				Table.Rows = append(Table.Rows, Temptable[index])
				max++
			}
			if max >= howmany {
				break
			}
		}
	}
	D.Table = append(D.Table, Table)
	*d = D
	if clear == true {
		Temptable = [][]string{}
	}
}

//ExportAsCSV lets you export temptable buffer as a CSV file.
func (d *Database) ExportAsCSV(filename string) {
	os.Remove(filename)
	csvFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	for index2 := range Temptable {
		csvwriter.Write(Temptable[index2])
	}
	if err := csvwriter.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	csvwriter.Flush()
	csvFile.Close()
}

//GrabTable lets you grab a table by name and then places the table columns/rows in the Temptable buffer.
//Passes to Buffer
func (d Database) GrabTable(tablename string) {
	if yes, index := d.verifytable(tablename); yes == true {
		Temptable = d.Table[index].Grab()
	}
}

//Normalize - normalises the data in a table.
//Affects table directly
func (d Database) Normalize(tablename string) {
	if yes, index := d.verifytable(tablename); yes == true {
		d.Table[index].Normalize()
	}
}

//Sum lets you count sum of total in a column of integers
//Passes to Buffer
func (d Database) Sum(tablename, columnname string) {
	if yes, index := d.verifytable(tablename); yes == true {
		Temptable = d.Table[index].Sum(columnname)
	}
}

//Clear deletes all data currently in temptable buffer
func (d Database) Clear() {
	Temptable = [][]string{}
}

//Shuffle lets you shuffle data in a table randomly
//Affects table directly
func (d *Database) Shuffle(tablename string) {
	D := *d
	if yes, tindex := d.verifytable(tablename); yes == true {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for index, i := range r.Perm(len(d.Table[tindex].Rows)) {
			D.Table[tindex].Rows[index] = D.Table[tindex].Rows[i]
		}
	}
	*d = D
}

//LoadDBase lets you load a database that you have previously saved
func (d *Database) LoadDBase(filename string) {
	item := *d
	jsonFile, _ := ioutil.ReadFile(filename + ".json")
	_ = json.Unmarshal([]byte(jsonFile), &item)
	*d = item
	fmt.Println("Loaded " + filename + "!")
}

//SaveDBase lets you save a database that you are currently working with
func (d Database) SaveDBase(filename string) {
	Base := &d
	output, err := json.MarshalIndent(Base, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = ioutil.WriteFile(filename+".json", output, 0755)
	fmt.Println("Saved " + filename + "!")
}

//Regex lets you grab a table where items within a column match a regular expression that you can pass in
//to the function and pulls whether they do or don't match (true/false)
//Passes to Buffer
func (d Database) Regex(tablename, columnname, regexquery string, boolean bool) {
	if yes, index := d.verifytable(tablename); yes == true {
		Temptable = d.Table[index].Regex(columnname, regexquery, boolean)
	}
}

//RegexReplace lets you replace substrings in strings with new string for rows that match a Regex.
//Affects table directly
func (d Database) RegexReplace(tablename, columnname, regexquery, oldstring, newstring string) {
	if yes, index := d.verifytable(tablename); yes == true {
		d.Table[index].RegexReplace(columnname, regexquery, oldstring, newstring)
	}
}

//CalculationInt lets you do calculations on the rows of a numerical column.
//You can pass a function of type func(column float64) float64.
//Affects table directly & adds another column (column name + _1D)
//Digits is how many decimal places the results will be
func (d Database) CalculationInt(tablename, columnname string, function table.Glfunc, decimals int) {
	if yes, index := d.verifytable(tablename); yes == true {
		d.Table[index].CalculationInt(columnname, function, decimals)
	}
}

//ConditionalInt lets you do conditionals on an integer column.
//You can pass a function of type func(x float64) bool.
//Passes to Buffer
//Results that are true will be returned
func (d Database) ConditionalInt(tablename, column string, function table.Cnfunc, decimals int) {
	if yes, index := d.verifytable(tablename); yes == true {
		Temptable = d.Table[index].ConditionalInt(column, function, decimals)
	}
}

//Calculation2DInt lets you do calculations on the rows of two numerical columns.
//You can pass a function of type func(x, y float64) float64.
//Affects table directly & adds another column (x name + _2D_ + y name)
//Digits is how many decimal places the results will be
func (d Database) Calculation2DInt(tablename, column1name, column2name string, function table.Gl2func, decimals int) {
	if yes, index := d.verifytable(tablename); yes == true {
		d.Table[index].Calculation2DInt(column1name, column2name, function, decimals)
	}
}

//Conditional2DInt lets you do conditionals on two integer columns at once.
//You can pass a function of type func(x, y float64) bool.
//Passes to Buffer
//Results that are true will be returned
func (d Database) Conditional2DInt(tablename, column1name, column2name string, function table.Cn2func, decimals int) {
	if yes, index := d.verifytable(tablename); yes == true {
		Temptable = d.Table[index].Conditional2DInt(column1name, column2name, function, decimals)
	}
}

//Conditional2DText lets you do conditionals on two string columns at once.
//You can pass a function of type func(x, y string) bool.
//Passes to Buffer
//Results that are true will be returned
func (d Database) Conditional2DText(tablename, column1name, column2name string, function table.Txt2func) {
	if yes, index := d.verifytable(tablename); yes == true {
		Temptable = d.Table[index].Conditional2DText(column1name, column2name, function)
	}
}

//DateToDays converts all dates in a column to days from today and then adds a new column called 'DateToDays'
//Directly Affects Table
//just choose table and column, parse is for what the date layout is i.e
//example - if it is SQL layout 2006-01-02T15:04:05-0700 then you want to write in '2006-01-02T15:04:05-0700'
func (d Database) DateToDays(tablename, columnname, parse string, delta float64) {
	if yes, index := d.verifytable(tablename); yes == true {
		d.Table[index].DateToDays(columnname, parse, delta)
	}
}

//Order re-orders a disorderly set of data by one column of integers
//Directly affects table.
//boolean true for ASC false for DESC
func (d Database) Order(tablename, columnname string, order bool) {
	if yes, index := d.verifytable(tablename); yes == true {
		d.Table[index].Order(columnname, order)
	}
}

//Plot - pass a table and a column to generate a plot of all fields against the column items.
//Max sample size is 155 - if one item has more than 155 samples it will display only a sample of the dataset.
func (d Database) Plot(table, namecolumn string) {
	if yes, index := d.verifytable(table); yes == true {
		d.Table[index].Plot(namecolumn)
	}
}

//AddIndex adds index column to wahtever is stored currently in Temptable
//Integer is where index starts from i.e 0 = starts from 0.
func (d Database) AddIndex(start int) {
	for index := range Temptable {
		if index == 0 {
			Temptable[index] = append(Temptable[index], "INDEX")
		} else {
			Temptable[index] = append(Temptable[index], strconv.Itoa(index+start-1))
		}
	}
}

//AddString adds a column containing string
//Passes to Buffer
func (d Database) AddString(columnname, addstring string) {
	for index := range Temptable {
		if index == 0 {
			Temptable[index] = append(Temptable[index], columnname)
		} else {
			Temptable[index] = append(Temptable[index], addstring)
		}
	}
}

//Row lets you grab a specific row number of items that have already been orderered.
//Passes to Buffer
func (d Database) Row(tablename, columnname string, rownumber int) {
	if yes, index := d.verifytable(tablename); yes == true {
		Temptable = d.Table[index].Row(columnname, rownumber)
	}
}

//Count lets you count the number of instances of all unique row items in column 'column' in table 'table'.
//Passes to Buffer
func (d Database) Count(tablename, columnname string) {
	if yes, index := d.verifytable(tablename); yes == true {
		Temptable = d.Table[index].Count(columnname)
	}
}

//Unpivot lets you unpivot data in table.
//Passes to Buffer
func (d Database) Unpivot(tablename, pivotcolumn string) {
	if yes, index := d.verifytable(tablename); yes == true {
		Temptable = d.Table[index].Unpivot(pivotcolumn)
	}
}

//Join is: left join table1 on table1.column1 = table2.column2
//Passes to Buffer
func (d Database) Join(table1, column1, table2, column2 string) {
	if yes, index := d.verifytable(table1); yes == true {
		if yes2, index2 := d.verifytable(table2); yes2 == true {
			Temptable = d.Table[index].Join(d.Table[index2], column1, column2)
		}
	}
}

//Timer lets you track how long queries take (try defer Timer(time.Now()))
func (d Database) Timer(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("This took %s", elapsed)
}

//Describe - pass a table & a dependent column to describe table around that column.
func (d Database) Describe(table, dependentcolumn string) {
	if yes, index := d.verifytable(table); yes == true {
		d.Table[index].Describe(dependentcolumn)
	}
}

//Columns returns a string array of the columns in a specific table 
//Used for when you want to do adjustments on each column via a loop
func (d Database) Columns(table string) []string {
	if yes, index := d.verifytable(table); yes == true {
		return d.Table[index].Columns
	}
	return nil
}

//Select lets you trim columns in temptable buffer to specific columns.
//You must pass in a 1D string array of column headers.
//Passes to Buffer
func (d Database) Select(columns []string) {
	replacers := []string{}
	for index := range columns {
		if output, item := grabsubstring(columns[index]); output != "" {
			replacers = append(replacers, output)
			columns[index] = item
		} else {
			replacers = append(replacers, item)
			columns[index] = item
		}
	}
	mainindex := []int{}
	newoutput := [][]string{}
	for tempindex := range Temptable[0] {
		for columnindex := range columns {
			if Temptable[0][tempindex] == columns[columnindex] {
				mainindex = append(mainindex, tempindex)
			}
		}
	}
	for index1 := range Temptable {
		output := []string{}
		if index1 != 0 {
			for index2 := range mainindex {
				output = append(output, Temptable[index1][mainindex[index2]])
			}
		} else {
			for indexx := range replacers {
				output = append(output, replacers[indexx])
			}

		}
		table.Addslice(&newoutput, output)
	}
	Temptable = newoutput
}

//Split - split a set of data up into two new tables (training / testing) at a certain ratio i.e 2 will be 50/50.
func (d *Database) Split(tablename, trainingname, testname string, ratio int) {
	D := *d
	Testtable := table.Table{}
	Testtable.Name = testname
	Trainingtable := table.Table{}
	Trainingtable.Name = trainingname
	if yes, tableindex := d.verifytable(tablename); yes == true {
		d.Shuffle(tablename)
		Testtable.Columns = d.Table[tableindex].Columns
		Trainingtable.Columns = d.Table[tableindex].Columns
		for index := range d.Table[tableindex].Rows {
			if index%ratio == 0 {
				Testtable.Rows = append(Testtable.Rows, d.Table[tableindex].Rows[index])
			} else {
				Trainingtable.Rows = append(Trainingtable.Rows, d.Table[tableindex].Rows[index])
			}
		}
	}
	D.Table = append(D.Table, Testtable)
	D.Table = append(D.Table, Trainingtable)
	fmt.Printf("Training size: %d\n", len(Trainingtable.Rows))
	fmt.Printf("Test size: %d\n", len(Testtable.Rows))
	*d = D
}

//KNNclass classifier using a training table to predict output on test table using identifier column.
//Passes to Buffer
func (d Database) KNNclass(trainingtable, testtable, identifiercolumn string, knumber int, trainingmode bool) {
	var trainingdata [][]float64
	var trainingname []string
	var testingdata [][]float64
	var testingname []string
	fmt.Println("KNN classification - importing data...")
	if yes, table1index := d.verifytable(trainingtable); yes == true {
		trainingname, trainingdata = d.Table[table1index].Grabdata(identifiercolumn)
	}
	if yes, table2index := d.verifytable(trainingtable); yes == true {
		testingname, testingdata = d.Table[table2index].Grabdata(identifiercolumn)
	}
	fmt.Println("KNN classification - processing...")
	Temptable = models.KNN(trainingdata, testingdata, trainingname, testingname, knumber, trainingmode, false)
}

//KNNreg regression using a training table to predict numerical output on test table using identifier column.
//Passes to Buffer
func (d Database) KNNreg(trainingtable, testtable, identifiercolumn string, knumber int) {
	var trainingdata [][]float64
	var trainingname []string
	var testingdata [][]float64
	var testingname []string
	fmt.Println("KNN regression - importing data...")
	if yes, table1index := d.verifytable(trainingtable); yes == true {
		trainingname, trainingdata = d.Table[table1index].Grabdata(identifiercolumn)
	}
	if yes, table2index := d.verifytable(trainingtable); yes == true {
		testingname, testingdata = d.Table[table2index].Grabdata(identifiercolumn)
	}
	fmt.Println("KNN regression - processing...")
	Temptable = models.KNN(trainingdata, testingdata, trainingname, testingname, knumber, false, true)
}

func outputcount(names []string) int {
	trainmap := make(map[string]int)
	iter := 0
	for index := range names {
		_, ok := trainmap[names[index]]
		if ok {
			continue
		} else {
			trainmap[names[index]] = iter
			iter++
		}
	}
	return iter
}

//NNtrain train a neural network using a training table.
//Need to provide training table, identifier column as well as number of hidden weights, epochs and learning rate.
//No need to provide number of inputs and outputs - this is calculated automatically to save troubles.
func (d *Database) NNtrain(trainingtable, identifiercolumn string, hidden, epochs int, learningrate float64) {
	var trainingdata [][]float64
	var trainingname []string
	fmt.Println("Importing data...")
	if yes, table1index := d.verifytable(trainingtable); yes == true {
		trainingname, trainingdata = d.Table[table1index].Grabdata(identifiercolumn)
		d.Net = models.CreateNN(len(trainingdata[0]), hidden, outputcount(trainingname), learningrate)
	}
	fmt.Println("Training NN...")
	d.Net.Train(trainingdata, trainingname, epochs)
}

//NNpredict use a trained neural network to predict another dataset, passing the table & identifier column.
func (d Database) NNpredict(trainingtable, identifiercolumn string) {
	var trainingdata [][]float64
	var trainingname []string
	fmt.Println("Importing data...")
	if yes, table1index := d.verifytable(trainingtable); yes == true {
		trainingname, trainingdata = d.Table[table1index].Grabdata(identifiercolumn)
	}
	fmt.Println("Predicting...")
	Temptable = d.Net.Predict(trainingdata, trainingname)
}
