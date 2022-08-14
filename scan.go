package main

import (
  "fmt"
  "os"
  "os/user"
  "strings"
  "log"
)

func scanGitFolders(folders []string, folder string)  []string{
  folder = strings.TrimSuffix(folder, "/")

  fd, err := os.Open(folder)
  if err != nil {
    log.Fatal(err)
  }

  files, err := fd.ReadDir(-1)
  fd.Close()
  if err != nil {
    log.Fatal(err)
  }

  var path string

  for _, file := range files {

    if file.IsDir() {
      path = folder + "/" + file.Name()
      if file.Name() == ".git" {
        path = strings.TrimSuffix(path, "/.git")
        fmt.Println(path)
        folders = append(folders)
        continue
      }
      if file.Name() == "node_modules" || file.Name() == "vendor" {
        continue
      }

      folders = scanGitFolders(folders, path)
    }
  }
  return folders
}

func recursiveScanPath(folder string) []string  {
  return scanGitFolders(make([]string, 0), folder)
}

func getFilePath() string {
  usr, err := user.Current()
  if err != nil {
    log.Fatal(err)
  }

  dotFile := usr.HomeDir + "/.git-going_local-stats"
  return dotFile
}

func parseLinesToSlice(filePath string) []string {
  
}

func addNewElementsToFile(filePath string, newRepos []string)  {
  existing := parseLinesToSlice(filePath)
  repos := joinSlices(newRepos, existing)
  dumpStringsToFile(repos, filePath)
}
// IDEA: Scans a given path for git repositories
func Scan(path string)  {

  fmt.Printf("Found the following folders:\n\n")
  repos := recursiveScanPath(path) // IDEA: get slice of strings of folder paths
  filePath := getFilePath() // IDEA: path of dotfile to write to
  addNewElementsToFile(filePath, repos) // IDEA: add slice contents to file
  fmt.Printf("\n\nSuccessfully Added\n\n")
}
