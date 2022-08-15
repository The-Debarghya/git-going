package main

import (
  "flag"
)

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
