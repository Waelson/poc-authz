# Serviço de Autorização baseado no SpiceDB (Google Zanzibar)
Esse projeto demonstra a integração entre uma aplicação web e o mecanismo de autorização SpiceDB desenvolvido pela Authzed baseado no paper Zanzibar da Google

## Pré requisites
- Golang 1.21
- Docker
## Instalação
Após assegur que o Docker está instalado em sua máquina. Exlsecute o comado abaixo para executar o SpiceDB em sua máquina. Atente-se para o parametro `--grpc-preshared-key` o valor dele será utilizado para conectar sua aplicação ao SpiceDB.

`
docker run --rm -p 50051:50051 authzed/spicedb serve --grpc-preshared-key "somerandomkeyhere"
`

## Realizando a autorização

Associa o usuário (sofia) a um recurso (10) com um relacionamento do tipo `reader`.

`
curl --location --request POST 'http://localhost:8080/authz/relationship' \
--header 'Content-Type: application/json' \
--data-raw '{
"resource": {
"namespace":"blog/post",
"id": "10"
},
"relation": "reader",
"subject": {
"namespace":"blog/user",
"id": "sofia"
}
}'
`

Checando a permissão

Verifica se o usuário `sofia` possui permissão de `read` no documento `10`

`
curl --location --request POST 'http://localhost:8080/authz/check-permission' \
--header 'Content-Type: application/json' \
--data-raw '{
"resource": {
"namespace":"blog/post",
"id": "10"
},
"permission": "read",
"subject": {
"namespace":"blog/user",
"id": "sofia"
}
}'
`