# Entity Status API

## Running the project

### Docker

```
$ docker build . -t entity-status-api:latest
$ docker run --rm --detach -p 8080:8080 entity-status-api:latest
$ cd app && go test -v 
```
