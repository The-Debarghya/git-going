# git-going

* A golang based command-line tool that helps to visualize your local git contributions

### Usage:

* If you have `golang` installed in your system, download and build it using:
```bash
git clone https://github.com/The-Debarghya/git-going
cd git-going/
go build .
```
 **OR:**
* Use `go get github.com/The-Debarghya/git-going` to download the package and install it using `go install`. <br>

**OR:**
* Download directly from the <a href="https://github.com/The-Debarghya/git-going/releases/tag/v1.0.0">Releases</a> section.

 **OPTIONS:**

* `git-going -add [path/to/your/folder]` : adds new entries of local git repositories and stores in *~/.git-going_local-stats*.

* `git-going -email email@yourdomain.com` : prints the contribution map of the stored folders.
