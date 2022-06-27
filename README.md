
# Domain Driven Design

> A News and Topic application

## Architecture

project layout applies Golang's Standard project layout
https://github.com/golang-standards/project-layout

## Build, run, automated tests

### build and run application


```shell script

make docker-start

```


### stop application

```shell script

make docker-stop

```


### automated tests

```shell script

make tests

```

# Usage :

endpoints URL is

- news http://localhost:8080/v1/news
- topic http://localhost:8080/v1/topic
- tag http://localhost:8080/v1/tag


### Create News

```shell script
curl -i -X POST http://localhost:8080/v1/news \
-H 'Content-Type: application/json' \
-d '{
	"title": "How to start investment",
	"writer": "setia budi",
	"content": "Lorem Ipsum is simply dummy text",
	"status": "publish",
	"tags": [1,2,3],
	"topic_id": 1
}'
```

### Create Topic

```shell script
curl -i -X POST http://localhost:8080/v1/topic \
-H 'Content-Type: application/json' \
-d '{"topic": "Investments"}'
```

### Create Tag

```shell script
curl -i -X POST http://localhost:8080/v1/tag \
-H 'Content-Type: application/json' \
-d '{"tag": "stock"}'
```

### Get all news with filter

```shell script
curl -i -X GET http://localhost:8080/v1/news/?status=publish&topic=1
```

### Get all news with filter date 

```shell script
curl -i -X GET http://localhost:8080/v1/news/?status=publish&topic=1&created_start=2022-06-20&created_end=2022-06-23
```
