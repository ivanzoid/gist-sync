
Use this script to maintain a set of github gists as files.

Usage:
 - install 'gist' program (https://github.com/defunkt/gist)
 - make sure you have Go (https://golang.org) installed
 - make a folder somewhere where you want to store your gists.
 - `go get github.com/ivanzoid/gist-sync ` or dl gist-sync.go and place it to this folder
 - place your gists as files to this folder
 - when you create a new or update existing gist file, just run: 'gist-sync gist.foo'
 - to sync all your gists, you may run 'gist-sync *.*'

