package main

import (
  "fmt"
  "os"
  "os/user"
  "strings"
  "log"
  "bufio"
  "io"
  "io/ioutil"
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

func openFile(filePath string) *os.File {
  f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0755)

  if err != nil {
    if os.IsNotExist(err) {
      // file doesn't exist
      _, err = os.Create(filePath)
      if err != nil {
        panic(err)
      }
    }else {
      // some other problem
      panic(err)
    }
  }
  // return the file descriptor
  return f
}

func parseLinesToSlice(filePath string) []string {
  f := openFile(filePath)
  defer f.Close() //postpone closing of the file descriptor

  var lines []string
  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }

  if err := scanner.Err(); err != nil{
    if err != io.EOF {
      panic(err)
    }
  }

  return lines
}

func sliceContains(slice []string, value string) bool {
  for _, v := range slice {
    if v == value {
      return true
    }
  }
  return false
}

func joinSlices(newRepos []string, existing []string) []string {
  for _, i := range newRepos {
    if !sliceContains(existing, i) {
      existing = append(existing, i)
    }
  }

  return existing
}

func dumpStringsToFile(repos []string, filePath string)  {
  content := strings.Join(repos, "\n")
  ioutil.WriteFile(filePath, []byte(content), 0755)
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
