# jsonbase
minimalistic declarative data manipulation library developed on 27th Sep.

The idea of this tool is that you can do very specific data manipulation in go.
I aim to make it as easy to use as possible & plan to extend on this and make it more useful over time.

Docs: https://pkg.go.dev/github.com/tbal999/jsonbase


What can you do with it?

## Import flat files, 1D or 2D slices.
- You can import csv files for example, or alternatively you could import a SQL query.

## Use Buffer
- Most of the queries get passed to a buffer table, which you print out via 'Print()'.
- In between transformations, you can save the buffer as a new table, and do further transformations. It's like a temperory table in SQL.
- You could join one table to another, save this output, then join the output to another table etc etc.

## Sum
- You can count sum of total in a column of integers.

## Count
- You can find the individual count of each unique item in a column.

## Regex
- You can use regular expressions on columns to find either matches or not-matches.

## Order 
- You can order specific columns of text or integers in either ascending or descending order.

## Row
- You can grab items (after ordering them) at a specific row. I.e the second oldest instance of each unique item.

## Unpivot
- You can unpivot data.

## Add Index
- You can add an index to the item.

## Left Join
- You can perform a left join on two tables - on one column match

## Replace strings
- You can use regex to find items in a column, and then replace them with new strings
- I.e find all items with unnecessary whitespace, and then delete that whitespace.

## Functions
- You can apply functions directly to integer columns and choose how many decimal places you want back.

## Conditionals
- You can apply conditional functions to integers in integer columns, and find either matches or not-matches.

## Date to days
- You can convert dates (in many different formats!) to days from delta.

## Column iteration
- You can grab the columns of a table and then pass a single column query through all columns.
- For example you could remove unnecessary whitespace from every single column.

## Machine learning / stats
- You can describe a dataset on a specific column via 'Describe()'
- You can plot scatterplots (making use of termui library)
- You can use KNN machine learning algorithm i built - bruteforce algorithm only so smaller datasets atm only.

Using a combination of the above, and the buffer, you can perform a lot of tasks.
