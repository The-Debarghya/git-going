package main

import (
  "time"

  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing/object"
)

// constants
const limit = 99999 // an upper limit
const totalDays = 183 //six months approximate

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

func printCommits(commits map[int]int)  {
  keys := sortMapToSlice(commits)
  cols := buildCols(keys, commits)
  printCells(cols)
}
// IDEA: generates a graph of git contributions
func Stats(email string)  {
  commits := processRepositories(email)
  printCommits(commits)
}
