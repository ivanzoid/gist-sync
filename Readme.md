
Use this script to maintain a set of github gists as files.

Usage:
 - install 'gist' program (https://github.com/defunkt/gist)
 - make sure you have Go (https://golang.org) installed
 - make a folder somewhere where you want to store your gists.
 - place gist-update.go (https://gist.github.com/ivanzoid/611177bbd3f5cb0604810f07080757b3#file-gist-update-go) to this folder
 - place your gists as files to this folder
 - when you create a new or update existing gist file, just run: './gist-update.go gist.foo'
 - to sync all your gists, you may run './gist-update.go *.*'

