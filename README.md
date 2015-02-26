lurker
======

Web crawler written in go designed to search and save comic data.

## Installation
	$ vagrant up
	$ vagrant ssh

Install dependencies

    $ cd /home/vagrant/go/src/github.com/comicgator/lurker
	$ go get ./...

Install cmd's

	$ cd cmd/lurker
	$ go install

Then you can run the basic lurker executable from anywhere

	$ lurker