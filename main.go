package main

import (
  "flag"
)

func main(){

  var parentFolder string
  var email string

  flag.StringVar(&parentFolder, "add", "", "Add a folder to scan for Git repositories")
  flag.StringVar(&email, "email", "example@yourdomain.com", "Email to scan")
  flag.Parse()

  if parentFolder != "" {
    Scan(parentFolder)
    return
  }

  Stats(email)
}
