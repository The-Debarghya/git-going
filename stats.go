package main

import (
  "fmt"
  "time"
  "sort"

  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing/object"
)

// globals
const limit = 99999 // an upper limit
const totalDays = 183 //six months approximate
const totalWeeks = 26 //six months approximate

type column []int
// IDEA: determines correct place of a commit in the commit map
// in order to fill the days which have no commits
func calcOffset() int {
  var offset int
  weekday := time.Now().Weekday()

  switch weekday {
  case time.Sunday:
    offset = 7
  case time.Monday:
    offset = 6
  case time.Tuesday:
    offset = 5
  case time.Wednesday:
    offset = 4
  case time.Thursday:
    offset = 3
  case time.Friday:
    offset = 2
  case time.Saturday:
    offset = 1
  }
  return offset

}

// IDEA: utility to get exact time of today's start i.e 00:00:00
func getBeginning(t time.Time) time.Time {
  year, month, day := t.Date()
  startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
  return startOfDay
}

func countDays(date time.Time) int {
  days := 0
  now := getBeginning(time.Now())
  for date.Before(now) {
    date = date.Add(time.Hour * 24)
    days++
    if days > totalDays {
      return limit
    }
  }

  return days
}

// IDEA: gets the commits in path and puts into the commit map
func fillCommits(email string, path string, commits map[int]int) map[int]int {
  repo, err := git.PlainOpen(path) // IDEA: instance of the git repo to be analyzed
  if err != nil {
    panic(err)
  }

  ref, err := repo.Head() // IDEA: gets the HEAD reference of the git repository
  if err != nil {
    panic(err)
  }

  iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()}) // IDEA: records all the commit logs
  if err != nil {
    panic(err)
  }
  // IDEA: iterate through each commit
  offset := calcOffset()
  err = iterator.ForEach(func(c *object.Commit) error {
    daysAgo := countDays(c.Author.When) + offset

    if c.Author.Email != email {
      return nil
    }

    if daysAgo != limit {
      commits[daysAgo]++
    }
    return nil

  })
  if err != nil {
    panic(err)
  }

  return commits
}

func processRepositories(email string) map[int]int {
  filePath := getFilePath()
  repos := parseLinesToSlice(filePath)
  daysInMap := totalDays

  commits := make(map[int]int, daysInMap)
  for i := daysInMap; i > 0; i-- {
    commits[i] = 0
  }

  for _, path := range repos {
    commits = fillCommits(email, path, commits)
  }

  return commits
}

// IDEA: takes a map and returns a slice with the map keys in order of values
func sortMapToSlice(m map[int]int) []int {
  var keys []int
  for k := range m {
    keys = append(keys, k)
  }

  sort.Ints(keys)
  return keys
}

// IDEA: generates a map with rows and columns ready to be printed to the screen
func buildCols(keys []int, commits map[int]int) map[int]column {

  cols := make(map[int]column)
  col := column{}

  for _, k := range keys {
    week := int(k/7)
    day := k % 7

    if day == 0 { // IDEA: If day is sunday create new column and fill it
      col = column{}
    }

    col = append(col, commits[k])
    if day == 6 { // IDEA: if day is saturday, add the week to cols map
      cols[week] = col
    }
  }

  return cols
}

func printMonths()  {
  week := getBeginning(time.Now()).Add(-(totalDays*time.Hour*24))
  month := week.Month()
  fmt.Printf("         ")
  for {
    if week.Month() != month {
      fmt.Printf("%s ", week.Month().String()[:3])
      month = week.Month()
    } else{
      fmt.Printf("    ")
    }

    week = week.Add(7*time.Hour*24)
    if week.After(time.Now()) { //surpasses the present date.
      break
    }
  }
  fmt.Printf("\n")
}

func printDayCol(day int) {
  out := "     "
  switch day {
  case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
  }

  fmt.Printf(out)
}

// IDEA: Prints each cell according to number of commits
func printCell(val int, this bool)  {

  escape := "\033[0;37;30m"
  switch {
    case val > 0 && val < 5:
        escape = "\033[1;30;47m"
    case val >= 5 && val < 10:
        escape = "\033[1;30;43m"
    case val >= 10:
        escape = "\033[1;30;42m"
  }

  if this {
    escape = "\033[1;37;45m"
  }

  if val == 0 {
    fmt.Printf(escape + "  - " + "\033[0m")
    return
  }

  str := "  %d "
  switch {
    case val >= 10:
        str = " %d "
    case val >= 100:
        str = "%d "
  }

  fmt.Printf(escape + str + "\033[0m", val)
}

func printCells(cols map[int]column)  {
  printMonths()
  for j := 6; j >= 0; j-- {
    for i := totalWeeks+1; i >= 0; i-- {
      if i == totalWeeks+1 {
        printDayCol(j)
      }

      if col, ok := cols[i]; ok {
        if i == 0 && j == calcOffset()-1 { //case for today
          printCell(col[j], true)
          continue
        } else {
          if len(col) > j {
            printCell(col[j], false)
            continue
          }
        }
      }

      printCell(0, false)
    }

    fmt.Printf("\n")
  }
}

func printCommits(commits map[int]int) {
  keys := sortMapToSlice(commits)
  cols := buildCols(keys, commits)
  printCells(cols)
}
// IDEA: generates a graph of git contributions
func Stats(email string)  {
  commits := processRepositories(email)
  printCommits(commits)
}
