dev:
	go build -o ./tmp ./cmd/api/main.go && air

tw:
	./tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/output.css --watch