# MBA Go Challenge 01

## Como rodar 

Primeiro roda o docker com o sqlite e depois o server com os comandos abaixo

Por usar gorm, certifique-se de ter instalado o compilador de C (gcc) no sistema operacional

```bash
go run server/main.go
```

e depois o cliente

```bash
go run client/main.go
```