
build:
	go build -o raw
install:
	mkdir -p /etc/raw
	cp ./raw /usr/bin/rawapp
	cp -n data/raw.service /etc/systemd/system/
.DEFAULT: build