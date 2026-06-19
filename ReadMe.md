# MBA Go Challenge 01

## Como rodar 

Por usar gorm, certifique-se de ter instalado o compilador de C (gcc) no sistema operacional

Depois so rodar o comando do server

```bash
go run server/main.go
```

e depois o cliente

```bash
go run client/main.go
```

se quiser ver os registros da tabela do sqlite (pacote sqlite3 instalado)

```
cd server
sqlite3 app.db
select * from dollar_quotation_dbs;
```