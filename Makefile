run:
	HOST=:3000 \
	ZAP_PROPERTIES_ENDPOINT=http://grupozap-code-challenge.s3-website-us-east-1.amazonaws.com/sources/source-2.json \
	go run main.go

make docker-build:
	docker build -t zap .

make docker-run:
	docker run --name zap -p 8080:8080 -d zap