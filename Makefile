build:
	docker build -t newforum .

run:build
	docker run --name forum --rm -p8080:8080 newforum