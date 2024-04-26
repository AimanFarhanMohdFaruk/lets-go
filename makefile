build:
	@go build -o bin/snippetbox cmd/web/*.go
	
run: build
	@bin/snippetbox

view:
	templ generate view