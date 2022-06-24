default: build

docker-start:
	docker-compose up --build -d

docker-stop:
	docker-compose down

tests :
	go test ./test/... -count=1 -p 1
