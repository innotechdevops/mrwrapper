# mrwrapper

MariaDB wrapper for Golang.

## Install

```shell
go get github.com/innotechdevops/mrwrapper
```

## How to use

- arguments

```go
args := []any{ "name": "x" }
```

- Count

```go
count := mrwraper.Count(conn, "SELECT COUNT(id) FROM X", args...)
```

- Select One

```go
data := mrwraper.SelectOne[Struct](conn, "SELECT * FROM X WHERE id = 1", args...)
```

- Select List

```go
data := mrwraper.SelectList[Struct](conn, "SELECT * FROM X", args...)
```

- Create

```go
id := 0
args := []any{
	"x",
}
tx, err := mrwrapper.Create(conn, "INSERT INTO X (name) VALUES (?)", []any{&id}, args...)
_ = tx.Commit()
```

- Update

```go
set := "name=:name"
params := map[string]interface{}{
	"id": 1, 
	"name": "x",
}
tx, err := mrwrapper.Update(conn, "UPDATE X SET %s WHERE id=:id", set, params)
_ = tx.Commit()
```

- Delete

```go
id := 1
tx, err := mrwrapper.Delete(conn, "DELETE FROM X WHERE id = ?", id)
_ = tx.Commit()
```

