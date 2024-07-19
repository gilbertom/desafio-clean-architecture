# Desafio Clean Architecture
Projeto do Desafio Clean Architecture para conclusão da Pós Graduação em Go Expert da Full Cycle.

<p align="center">
  <img src="https://blog.golang.org/gopher/gopher.png" alt="">
</p>

Esta é uma aplicação desenvolvida em Go, cujo objetivo é realizar a criação de Orders e consultar todas as Orders através de requisições feitas através de HTTP, gRPC e GraphQL, tudo numa única aplicação.

<br>

## Principais Recursos Utilizados

- **Go**
- **Go-Chi**
- **Wire**
- **grpc**
- **gqlgen**
- **Viper**
- **MySQL**
- **RabbitqMQ**

<br>
## Índice

- [Instalação](#instalação)
- [Como Usar](#como-usar)
- [Portas por serviço](#portas-por-serviço)
- [Contato](#contato)
- [Agradecimentos](#agradecimentos)

<br>

## Instalação

```sh
$ git clone https://github.com/gilbertom/desafio-clean-architecture.git
$ docker-compose up --build
```
<br>

## Como Usar

Exemplo de uso utilizando chamada HTTP:

No diretório /api temos dois arquivos. Um chamado 'create_order.http' e outro 'query_order.http'. 

Abra no Visual Studio Code o arquivo 'create_order.http'. Garanta que a Extensão 'HTTP Client' esteja instalada, realize as alterações desejadas no body e clique em 'Send Request'.

Request
```sh
POST http://localhost:8000/order HTTP/1.1
Host: localhost:8000
Content-Type: application/json

{
    "id":"Order 1",
    "price": 100.5,
    "tax": 0.5
}
```

Response

  Em caso de sucesso:

  ```sh
  HTTP/1.1 200 OK
  Date: Thu, 18 Jul 2024 23:16:02 GMT
  Content-Length: 59
  Content-Type: text/plain; charset=utf-8
  Connection: close

  {
    "id": "Order 1",
    "price": 100.5,
    "tax": 0.5,
    "final_price": 101
  }
  ```

  Em caso de insucesso:

  ```sh
  HTTP/1.1 500 Internal Server Error
  Content-Type: text/plain; charset=utf-8
  X-Content-Type-Options: nosniff
  Date: Thu, 18 Jul 2024 23:18:28 GMT
  Content-Length: 71
  Connection: close

  Error 1062 (23000): Duplicate entry 'Order 1' for key 'orders.PRIMARY'
  ```

Exemplo de uso utilizando gRPC:

``` sh
$ evans -r repl


  ______
 |  ____|
 | |__    __   __   __ _   _ __    ___
 |  __|   \ \ / /  / _. | | '_ \  / __|
 | |____   \ V /  | (_| | | | | | \__ \
 |______|   \_/    \__,_| |_| |_| |___/

 more expressive universal gRPC client


pb.OrderService@127.0.0.1:50051> call CreateOrder
id (TYPE_STRING) => Order RPC 1
price (TYPE_FLOAT) => 10001
tax (TYPE_FLOAT) => 101.4
{
  "finalPrice": 10102.4,
  "id": "Order RPC 1",
  "price": 10001,
  "tax": 101.4
}
```
<br>

Exemplo de Consulta de Orders usando gRPC:

```sh
pb.OrderService@127.0.0.1:50051> call ListOrder
{
  "orders": [
    {
      "finalPrice": 101,
      "id": "Order 1",
      "price": 100.5,
      "tax": 0.5
    },
    {
      "finalPrice": 10102.4,
      "id": "Order RPC 1",
      "price": 10001,
      "tax": 101.4
    }
  ]
}
```

<br>

Exemplo de uso utilizando GraphQL para criar uma Order:

Acesse o endereço http://localhost:8080/ no seu browser.

Request

```sh
mutation createOrder {
  createOrder(input: {id: "Order GraphQL", Price: 1000.1, Tax:10.9}) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

Response

```sh
{
  "data": {
    "createOrder": {
      "id": "Order GraphQL",
      "Price": 1000.1,
      "Tax": 10.9,
      "FinalPrice": 1011
    }
  }
}
```

<br>

## Portas por serviço

**Serviço HTTP** - Responde na porta 8000  
**Serviço gRPC** - Responde na porta 50051  
**Serviço GraphQL** - Responde na porta 8080  

<br>

## Contato
Para entrar em contato com o desenvolvedor deste projeto:
[gilbertomakiyama@gmail.com](mailto:gilbertomakiyama@gmail.com)

<br>

## Agradecimentos
Gostaria de expressar minha sincera gratidão a todo o time do curso de Pós-Graduação em Go Avançado da FullCycle pelo empenho, dedicação e excelência no ensino. Suas contribuições foram fundamentais para o meu desenvolvimento e sucesso. Muito obrigado!