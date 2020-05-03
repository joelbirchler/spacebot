HOST=gene

all: build upload run

build:
	GOARM=6 GOARCH=arm GOOS=linux go build -o ./bin/spacebot

upload:
	scp ./bin/spacebot pi@${HOST}:~/spacebot

run:
	ssh -t pi@${HOST} "sudo ~/spacebot"