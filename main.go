package main

import (
  "flag"
  "fmt"
)

// IDEA: Scans a given path for git repositories
func Scan(path string)  {

  fmt.Printf("Found the following folders:\n\n")
  repos := recursiveScanPath(path) // IDEA: get slice of strings of folder paths
  filePath := getFilePath() // IDEA: path of dotfile to write to
  addNewElementsToFile(filePath, repos) // IDEA: add slice contents to file
  fmt.Printf("\n\nSuccessfully Added\n\n")
}

// IDEA: generates a graph of git contributions
func Stats(email string)  {
  print("stats")
}

func main(){

  var folder string
  var email string

  flag.StringVar(&folder, "add", "", "Add a folder to scan for Git repositories")
  flag.StringVar(&email, "email", "example@yourdomain.com", "Email to scan")
  flag.Parse()

  if folder != "" {
    Scan(folder)
    return
  }

  Stats(email)
}
