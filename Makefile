# Made by Leif Terje Fonnes, authumn 2017 for TV2
# requires that bump version is installed	go get github.com/Shyp/bump_version


.PHONY: test

bump:
	bump_version patch amqphook.go 

push: 
	git push origin --tags

release: test bump push

test:
	go test ./... -v
