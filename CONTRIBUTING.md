# Contributing

Send it!

### Dependencies

```bash
sudo apt install -y sqlite3 golang make
```

### Run

```bash
make dev
```

### Database Migrations

Create new migration with the current timestamp

```bash
touch sql/$(date +%Y%m%d%H%M%S).sql
```
