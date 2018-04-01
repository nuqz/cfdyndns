CF_API_EMAIL=
CF_API_KEY=

.PHONY: binary

binary:
	go build

local: binary
	CF_API_EMAIL=$(CF_API_EMAIL) \
	CF_API_KEY=$(CF_API_KEY) \
	./cfdyndns -interval=5s

install_service:
	sudo cp ./cfdyndns.service /etc/systemd/system/
	sudo systemctl enable cfdyndns.service
