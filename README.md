# TODO

- proposals
- db backup cronjob
- reverse proxy

### Dev

```bash
go run cmd/main.go
```

### Deploy

```bash
docker build -t rishunnik .
docker run -p 8080:8080 --name rishunnik --rm rishunnik
```
