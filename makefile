build:
	@templ generate view
	@go build -o bin/snippetbox cmd/web/*.go
	
run: build
	@bin/snippetbox

view:
	templ generate view

templ:
	@templ generate --watch -proxy=http://localhost:3000