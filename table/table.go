package table

import (
	"crypto/sha256"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"math/rand"

	ui "github.com/gizak/termui/v3"              //for scatterplots
	widgets "github.com/gizak/termui/v3/widgets" //for scatterplots
)

//Table is a struct containing the name, columns and rows of a table.
type Table struct {
	Name    string
	Columns []string
	Rows    [][]string
}

var s sync.Mutex //used in some situations i.e hashing.

type Glfunc func(float64) float64

type Gl2func func(float64, float64) float64

type Cnfunc func(float64) bool

type Cn2func func(float64, float64) bool

type Txt2func func(string, string) bool

///// HELPER FUNCTIONS /////

func datetodays(x, parse string, delta float64) string {
	t, _ := time.Parse(parse, x)
	duration := time.Now().Add(time.Duration(delta) * (time.Hour * 24)).Sub(t)
	y := fmt.Sprintf("%d", int(duration.Hours()/24))
	return y
}

func join2D(x, y []string) []string {
	z := append(x, y...)
	return z
}

func isletter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func put1D(hash1dtable *map[[32]byte][]string, item []string) {
	h := *hash1dtable
	s.Lock()
	defer s.Unlock()
	i := hash1D(item)
	h[i] = item
	*hash1dtable = h
}

func Addslice(two *[][]string, one []string) {
	a := *two
	a = append(a, one)
	*two = a
}

func hash(k string) [32]byte {
	h := sha256.Sum256([]byte(k))
	return h
}

func hash1D(k []string) [32]byte {
	join := strings.Join(k, "")
	h := sha256.Sum256([]byte(join))
	return h
}

func get1D(hash1Dtable map[[32]byte][]string, item []string) []string {
	s.Lock()
	defer s.Unlock()
	i := hash1D(item)
	return hash1Dtable[i]
}

func contains(arr []string, str string, howmany int) (int, bool) {
	x := 0
	for index := range arr {
		if arr[index] == str {
			x++
		}
	}
	if x != 0 {
		return x, true
	}
	return 0, false
}

func remove1D(slice []string, s int) []string {
	if s != len(slice) {
		slice = append(slice[:s], slice[s+1:]...)
		return slice
	} else if s == len(slice) {
		slice = slice[:len(slice)-1]
		return slice
	}
	return slice
}

func samplerate(x []float64) []float64 {
	sample := []float64{}
	var lowerquartile float64
	var upperquartile float64
	lNumberf :=  math.Floor(0.25 * float64(len(x)))
	lNumber := int(lNumberf)
	uNumberf :=  math.Floor(0.75 * float64(len(x)))
	uNumber := int(uNumberf)
	if len(x) % 2 == 0 {
		upperquartile =  (x[uNumber-1] + x[uNumber]) / 2
		lowerquartile = (x[lNumber-1] + x[lNumber]) / 2
	} else {
		upperquartile = x[uNumber]
		lowerquartile =  x[lNumber]
	}
	interquartile := upperquartile-lowerquartile
	min := lowerquartile - 1.5*interquartile
	max := upperquartile + 1.5*interquartile
	r := rand.New(rand.NewSource(time.Now().Unix()))
	counter := 0
	for _, i := range r.Perm(len(x)) {
		if counter > 155 {
			break
		}
		if x[i] >= min && x[i] <= max {
			sample = append(sample, x[i])
			counter++
		}
	}
	return sample
}

func scatterinit(x int, mainlist []map[string][]float64, columns []string, column string) {
	p2 := widgets.NewPlot()
	p3 := widgets.NewParagraph()
	p2.Marker = widgets.MarkerDot
	output := mainlist[x]
	p2.Data = make([][]float64, len(output))
	var i = 0
	var text = "Key: \n"
	length := 0
	for index, plotdata := range output {
		p2.Data[i] = plotdata
		if len(plotdata)+1 > 155 {
			p2.Data[i] = samplerate(p2.Data[i])
		}
		sort.Float64s(p2.Data[i])
		text += index + ": " + color(i+1) + "\n"
		if i != len(output)-1 {
		} else {
			text += `
			
			Instructions: 
			press q to quit
			press a and d to
			change field`
		}
		length = len(p2.Data[i])
		i++
	}
	p2.SetRect(0, 0, length, 30)
	p3.SetRect(length, 0, length+20, 30)
	p2.AxesColor = ui.ColorWhite
	p2.PlotType = widgets.ScatterPlot
	p2.Title = columns[x] + " for different " + column
	p3.Text = text
	p3.Border = false
	p3.TextStyle.Fg = ui.ColorBlue
	ui.Render(p2, p3)
}

func collate(main string, row, columns []string) [][]string {
	output := [][]string{}
	rowindex := 0
	newcolumns := []string{}
	newrow := []string{}
	var first bool
loop0:
	for index := range columns {
		if first == false {
			if columns[index] == main {
			newcolumns = append(newcolumns, columns[index])
			newrow = append(newrow, row[index])
			first = true
			goto loop0
			}
		} else {
			if columns[index] != main {
			newcolumns = append(newcolumns, columns[index])
			newrow = append(newrow, row[index])
			}
		}
	}
loop:
	columnindex := 0
	for columnindex < len(newcolumns) {
		item := []string{}
		if newcolumns[columnindex] == main {
			item = append(item, newrow[columnindex])
			if rowindex < len(newrow)-1 {
				rowindex++
				item = append(item, newcolumns[rowindex])
			} else {
				break
			}
		}
		if len(newcolumns) > 1 {
			item = append(item, newrow[rowindex])
			output = append(output, item)
			goto loop
		}
	}
	return output
}

func color(index int) string {
	switch index {
	case 0:
		return "black"
	case 1:
		return "red"
	case 2:
		return "green"
	case 3:
		return "yellow"
	case 4:
		return "blue"
	case 5:
		return "magenta"
	case 6:
		return "cyan"
	case 7:
		return "white"
	}
	return ""
}

func get(table map[[32]byte]string, item string) (string, [32]byte) {
	s.Lock()
	defer s.Unlock()
	i := hash(item)
	return table[i], i
}

func put(table *map[[32]byte]string, item string) {
	t := *table
	s.Lock()
	defer s.Unlock()
	i := hash(item)
	t[i] = item
	*table = t
}

func hashadd(table *map[[32]byte]int, item string) {
	t := *table
	s.Lock()
	defer s.Unlock()
	i := hash(item)
	if t[i] == 0 {
		t[i] = 1
	} else {
		t[i] = t[i] + 1
	}
}

func hashget(table map[[32]byte]int, item string) int {
	s.Lock()
	defer s.Unlock()
	i := hash(item)
	return table[i]
}

func joincolumns(table1, table2 Table) []string {
	output := []string{}
	for index := range table1.Columns {
		output = append(output, table1.Name+"_"+table1.Columns[index])
	}
	for index := range table2.Columns {
		output = append(output, table2.Name+"_"+table2.Columns[index])
	}
	return output
}

///// HELPER FUNCTIONS /////

func (t Table) mean(columnindex, dependentcolumn int, item string) (float64, float64) {
	var sum float64
	var howmany float64
	for index := range t.Rows {
		if t.Rows[index][dependentcolumn] == item {
			howmany++
			x, _ := strconv.ParseFloat(t.Rows[index][columnindex], 64)
			sum += x
		}
	}
	mean := sum / howmany
	var variance float64
	for index := range t.Rows {
		if t.Rows[index][dependentcolumn] == item {
			x, _ := strconv.ParseFloat(t.Rows[index][columnindex], 64)
			variance += (x - mean) * (x - mean)
		}
	}
	return mean, math.Sqrt(variance) //standard deviation
}

func (t Table) Grab() [][]string {
	output := [][]string{}
	output = append(output, t.Columns)
	for index := range t.Rows {
		output = append(output, t.Rows[index])
	}
	return output
}

func (t Table) plotgrab(identifiercolumn, datacolumn string) (names []string, data [][]float64) {
	column1index, real := t.verifycolumn(identifiercolumn)
	column2index, real2 := t.verifycolumn(datacolumn)
	if real == true && real2 == true {
		for index := range t.Rows {
			item := []float64{}
			for columnindex := range t.Rows[index] {
				if columnindex == column1index {
					names = append(names, t.Rows[index][column1index])
				}
				if columnindex == column2index {
					x, _ := strconv.ParseFloat(t.Rows[index][columnindex], 64)
					item = append(item, x)
				}
			}
			data = append(data, item)
		}
	}
	return names, data
}

func (t Table) collectdata(column, column2 string) map[string][]float64 {
output := make(map[string][]float64)
	_, real := t.verifycolumn(column)
	if real == true {
		names, data := t.plotgrab(column, column2)
		for index := range data {
			for i := range data[index] {
				output[names[index]] = append(output[names[index]], data[index][i])
			}
		}
	}
	return output
}

func (t Table) meanall(column int) (float64, float64) {
	var sum float64
	var howmany = float64(len(t.Rows))
	for index := range t.Rows {
		x, _ := strconv.ParseFloat(t.Rows[index][column], 64)
		sum += x
	}
	mean := sum / howmany
	var variance float64
	for index := range t.Rows {
		x, _ := strconv.ParseFloat(t.Rows[index][column], 64)
		variance += (x - mean) * (x - mean)
	}
	return mean, math.Sqrt(variance) //standard deviation
}

func (t *Table) Normalize() {
	T := *t
	for columnindex := range T.Columns {
		mean, dev := t.meanall(columnindex)
		for row := range T.Rows {
			if isletter(T.Rows[row][columnindex]) == false {
				x, _ := strconv.ParseFloat(T.Rows[row][columnindex], 64)
				x = (x - mean) / dev
				T.Rows[row][columnindex] = strconv.FormatFloat(x, 'f', 2, 64)
			}
		}
	}
	*t = T
}

func (t Table) verifycolumn(Column string) (int, bool) {
	for index := range t.Columns {
		if t.Columns[index] == Column {
			return index, true
		}
	}
	fmt.Println("Column " + Column + " does not exist in table " + t.Name)
	return 0, false
}

func (t Table) Sum(column string) [][]string {
	var sum float64
	output := [][]string{}
	output = append(output, t.Columns)
	columnindex, real := t.verifycolumn(column)
	if real == true {
		for index2 := range t.Rows {
			number, err := strconv.ParseFloat(t.Rows[index2][columnindex], 64)
			if err != nil {
				fmt.Println("There are non-integers in this column. Cannot sum.")
				return nil
			}
			sum += number
		}
		sumstring := []string{}
		sumstring = append(sumstring, strconv.FormatFloat(sum, 'f', 2, 64))
		output = append(output, sumstring)
		return output
	}
	return nil
}

func (t Table) Regex(column, regexquery string, boolean bool) [][]string {
	output := [][]string{}
	output = append(output, t.Columns)
	columnindex, real := t.verifycolumn(column)
	if real == true {
		for index := range t.Rows {
			if match, _ := regexp.MatchString(regexquery, t.Rows[index][columnindex]); match == boolean {
				output = append(output, t.Rows[index])
			}
		}
		return output
	}
	return nil
}

func (t *Table) RegexReplace(column, regexquery, oldstring, newstring string) {
	T := *t
	output := [][]string{}
	output = append(output, t.Columns)
	columnindex, real := t.verifycolumn(column)
	if real == true {
		for index := range t.Rows {
			if match, _ := regexp.MatchString(regexquery, t.Rows[index][columnindex]); match == true {
				T.Rows[index][columnindex] = strings.Replace(T.Rows[index][columnindex], oldstring, newstring, -1)
			}
		}
	}
	*t = T
}

func (t *Table) CalculationInt(columnname string, function Glfunc, decimals int) {
	T := *t
	columnindex, real := t.verifycolumn(columnname)
	if real == true {
		t.Columns = append(t.Columns, columnname+"_1D")
		for index := range t.Rows {
			number, _ := strconv.ParseFloat(T.Rows[index][columnindex], 64)
			T.Rows[index][columnindex] = strconv.FormatFloat(function(number), 'f', decimals, 64)
		}
	}
	*t = T
}

func (t *Table) Calculation2DInt(column1, column2 string, function Gl2func, decimals int) {
	T := *t
	column1index, real1 := t.verifycolumn(column1)
	column2index, real2 := t.verifycolumn(column2)
	if real1 == true && real2 == true {
		T.Columns = append(T.Columns, column1+"_2D_"+column2)
		for index := range t.Rows {
			number1, _ := strconv.ParseFloat(T.Rows[index][column1index], 64)
			number2, _ := strconv.ParseFloat(T.Rows[index][column2index], 64)
			T.Rows[index] = append(T.Rows[index], strconv.FormatFloat(function(number1, number2), 'f', decimals, 64))
		}
		*t = T
	}
}

func (t Table) ConditionalInt(columnname string, function Cnfunc, decimals int) [][]string {
	output := [][]string{}
	columnindex, real := t.verifycolumn(columnname)
	if real == true {
		output = append(output, t.Columns)
		for index := range t.Rows {
			number, _ := strconv.ParseFloat(t.Rows[index][columnindex], 64)
			if function(number) == true {
				output = append(output, t.Rows[index])
			}
		}
		return output
	}
	return nil
}

func (t Table) Conditional2DInt(column1, column2 string, function Cn2func, decimals int) [][]string {
	output := [][]string{}
	column1index, real1 := t.verifycolumn(column1)
	column2index, real2 := t.verifycolumn(column2)
	if real1 == true && real2 == true {
		output = append(output, t.Columns)
		for index := range t.Rows {
			number1, _ := strconv.ParseFloat(t.Rows[index][column1index], 64)
			number2, _ := strconv.ParseFloat(t.Rows[index][column2index], 64)
			if function(number1, number2) == true {
				output = append(output, t.Rows[index])
			}
		}
		return output
	}
	return nil
}

func (t Table) Conditional2DText(column1, column2 string, function Txt2func) [][]string {
	output := [][]string{}
	column1index, real1 := t.verifycolumn(column1)
	column2index, real2 := t.verifycolumn(column2)
	if real1 == true && real2 == true {
		output = append(output, t.Columns)
		for index := range t.Rows {
			string1 := t.Rows[index][column1index]
			string2 := t.Rows[index][column2index]
			if function(string1, string2) == true {
				output = append(output, t.Rows[index])
			}
		}
		return output
	}
	return nil
}

func (t Table) DateToDays(column, parse string, delta float64) {
	columnindex, real := t.verifycolumn(column)
	if real == true {
		t.Columns = append(t.Columns, column+"_days")
		for index := range t.Rows {
			t.Rows[index] = append(t.Rows[index], datetodays(t.Rows[index][columnindex], parse, delta))
		}
	}
}

func (t Table) Order(column string, order bool) {
	columnindex, real := t.verifycolumn(column)
	if real == true {
		if order == true {
			sort.SliceStable(t.Rows, func(i, j int) bool {
				itext := isletter(t.Rows[i][columnindex])
				jtext := isletter(t.Rows[j][columnindex])
				if itext && jtext == true {
					return t.Rows[i][columnindex] < t.Rows[j][columnindex]
				}
				inumber, _ := strconv.ParseFloat(t.Rows[i][columnindex], 64)
				jnumber, _ := strconv.ParseFloat(t.Rows[j][columnindex], 64)
				return inumber < jnumber
			})
		} else {
			sort.SliceStable(t.Rows, func(i, j int) bool {
				itext := isletter(t.Rows[i][columnindex])
				jtext := isletter(t.Rows[j][columnindex])
				if itext && jtext == true {
					return t.Rows[i][columnindex] > t.Rows[j][columnindex]
				}
				inumber, _ := strconv.ParseFloat(t.Rows[i][columnindex], 64)
				jnumber, _ := strconv.ParseFloat(t.Rows[j][columnindex], 64)
				return inumber > jnumber
			})
		}
	}
}

func (t Table) Row(column string, rownumber int) [][]string {
	columnindex, real := t.verifycolumn(column)
	if real == true {
		stringlist := []string{}
		output := [][]string{}
		output = append(output, t.Columns)
		for index := range t.Rows {
		inter:
			stringlist = append(stringlist, t.Rows[index][columnindex])
			count, yes := contains(stringlist, t.Rows[index][columnindex], rownumber)
			if rownumber == count && yes == true {
				output = append(output, t.Rows[index])
				goto inter
			}
		}
		return output
	}
	return nil
}

func (t *Table) Unpivot(maincolumn string) [][]string {
	newcolumns := []string{maincolumn, "UNPIVOT", "ITEMS"}
	item := [][][]string{}
	final := [][]string{}
	final = append(final, newcolumns)
	for rowindex := range t.Rows {
		in := t.Rows[rowindex]
		output := collate(maincolumn, in, t.Columns)
		item = append(item, output)
	}
	for index1 := range item {
		for index2 := range item[index1] {
			final = append(final, item[index1][index2])
		}
	}
	return final
}

func (t Table) Count(column string) [][]string {
	output := [][]string{}
	hashtable := make(map[[32]byte]string)
	hashcount := make(map[[32]byte]int)
	header := []string{t.Name + "_" + column + "_ITEM", t.Name + "_" + column + "_COUNT"}
	output = append(output, header)
	columnindex, real := t.verifycolumn(column)
	if real == true {
		for index := range t.Rows {
			if grab, _ := get(hashtable, t.Rows[index][columnindex]); grab == "" {
				put(&hashtable, t.Rows[index][columnindex])
				hashadd(&hashcount, t.Rows[index][columnindex])
			} else {
				hashadd(&hashcount, t.Rows[index][columnindex])
			}
		}
	}
	for key := range hashtable {
		str := hashtable[key] + "," + strconv.Itoa(hashget(hashcount, hashtable[key]))
		newitem := strings.Split(str, ",")
		Addslice(&output, newitem)
	}
	return output
}

func (t Table) Join(t2 Table, tcolumn, t2column string) [][]string {
	column1index, real1 := t.verifycolumn(tcolumn)
	column2index, real2 := t.verifycolumn(t2column)
	if real1 && real2 == true {
		hash1Dtable := make(map[[32]byte][]string)
		output := [][]string{}
		output = append(output, joincolumns(t, t2))
		var abool bool
		for yindex := range t.Rows {
			abool = false
			for xindex := range t2.Rows {
				if t.Rows[yindex][column1index] == t2.Rows[xindex][column2index] {
					abool = true
					if item := get1D(hash1Dtable, join2D(t.Rows[yindex], t2.Rows[xindex])); item == nil {
						put1D(&hash1Dtable, join2D(t.Rows[yindex], t2.Rows[xindex]))
						output = append(output, join2D(t.Rows[yindex], t2.Rows[xindex]))
					}
				}
			}
			if abool == false {
				var null = []string{"<-NULL->"}
				output = append(output, join2D(t.Rows[yindex], null))
			}
		}
		return output
	}
	return nil
}

func (t Table) Describe(column string) {
	column1index, real := t.verifycolumn(column)
	if real == true {
		columnnames := []string{}
		for index2 := range t.Rows {
			if _, yes := contains(columnnames, t.Rows[index2][column1index], 1); yes == false {
				columnnames = append(columnnames, t.Rows[index2][column1index])
			}
		}
		for index := range columnnames {
			fmt.Println("Describe for " + columnnames[index])
			for index2 := range t.Columns {
				if t.Columns[index2] != column {
					mean, dev := t.mean(index2, column1index, columnnames[index])
					fmt.Printf("For column '%s' - average is %f, standard deviation is %f\n", t.Columns[index2], mean, dev)
					
				}
			}
			fmt.Println("")
		}
		for index2 := range t.Columns {
				t.Focusdata(column,t.Columns[index2]) 
		}
	}
	fmt.Println("Consider removing outliers from source data")
}

func (t Table) Grabdata(identifiercolumn string) (names []string, data [][]float64) {
	column1index, real := t.verifycolumn(identifiercolumn)
	if real == true {
		for index := range t.Rows {
			item := []float64{}
			for columnindex := range t.Rows[index] {
				if columnindex == column1index {
					names = append(names, t.Rows[index][column1index])
				} else {
					x, _ := strconv.ParseFloat(t.Rows[index][columnindex], 64)
					item = append(item, x)
				}
			}
			data = append(data, item)
		}
	}
	return names, data
}

func (t Table) Focusdata(column1, column2 string ) {
	output := t.collectdata(column1,column2)
	var item string
	list := []float64{}
	for index, plotdata := range output {
		item = index
		list = plotdata
		fmt.Println("For: "+item+" - column: " +column2)
		sort.Float64s(list)
		minoutliers, maxoutliers := t.stats(list)
		fmt.Printf("Min outliers: %f\n", minoutliers)
		fmt.Printf("Max outliers: %f\n", maxoutliers)
		fmt.Println("")
	}
}

func (t Table) stats(x []float64) ([]float64, []float64) {
	minout := []float64{}
	maxout := []float64{}
	var lowerquartile float64
	var upperquartile float64
	lNumberf :=  math.Floor(0.25 * float64(len(x)))
	lNumber := int(lNumberf)
	uNumberf :=  math.Floor(0.75 * float64(len(x)))
	uNumber := int(uNumberf)
	if len(x) % 2 == 0 {
		upperquartile =  (x[uNumber-1] + x[uNumber]) / 2
		lowerquartile = (x[lNumber-1] + x[lNumber]) / 2
	} else {
		upperquartile = x[uNumber]
		lowerquartile =  x[lNumber]
	}
	interquartile := upperquartile-lowerquartile
	min := lowerquartile - 1.5*interquartile
	max := upperquartile + 1.5*interquartile
	for index := range x {
		if x[index] < min {
			minout = append(minout, x[index])
		}
		if x[index] > max {
			maxout = append(maxout, x[index])
		}
	}
	return minout, maxout
}

func (t Table) Plot(column string) {
	mainlist := []map[string][]float64{}
	columns := []string{}
	for index := range t.Columns {
		if t.Columns[index] != column {
			output := t.collectdata(column,t.Columns[index])
			columns = append(columns, t.Columns[index])
			mainlist = append(mainlist, output)
		}
	}
	if err := ui.Init(); err != nil {
		fmt.Println("Error generating plot")
		return
	}
	x := len(mainlist)-1
	defer ui.Close()
	scatterinit(1, mainlist, columns, column)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "a":
			if x == 0 {
				x = len(mainlist)-1
			} else {
				x--
			}
			scatterinit(x, mainlist, columns, column)
		case "d":
			if x == len(mainlist)-1 {
				x = 0
			} else {
				x++
			}
			scatterinit(x, mainlist, columns, column)
		}

	}
}
