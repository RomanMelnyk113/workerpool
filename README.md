# Worker pool

Worker pool implementation where it can accept any function to be executed. User can define number of workers and add
limitation for number of tasks to execute
Note: currently it is simply taking testing file https://s3.amazonaws.com/alexa-static/top-1m.csv.zip and sending 
requests to all URLs in the list based on limits

## Usage
Build the project first
```sh
go build -o workerpool ./cmd/workerpool/main.go
```

then run
``` sh 
./workerpool -workers=5 -tasks=100
```
