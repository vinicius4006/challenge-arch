# Como Rodar o Projeto

### Pré-requisitos:
- **Docker**: Certifique-se de que o Docker está instalado na máquina.

### Instruções de Execução:
1. Após garantir que o Docker está instalado, rode o seguinte comando para construir e iniciar os serviços:
   ```bash
   docker compose up --build
   ```

2. Aguarde até que todos os serviços estejam ativos e prontos.

---

## Testes Via REST

### 1. Criar Pedido:
Execute o seguinte comando `curl` para criar um pedido:

```bash
curl -X POST http://localhost:8000/order \
-H "Content-Type: application/json" \
-d '{
    "id": "tame",
    "price": 100.5,
    "tax": 0.5
}'
```

### 2. Listar Pedidos:
Execute este comando para listar todos os pedidos:

```bash
curl -X GET http://localhost:8000/list \
-H "Content-Type: application/json"
```

---

## Testes Via GraphQL

### Acessar GraphQL:
1. Acesse o serviço GraphQL em: [http://localhost:8080](http://localhost:8080)

### 1. Criar Pedido:
Utilize a seguinte **mutation** para criar um pedido:

```graphql
mutation CreateOrder {
  CreateOrder(input: {id: "teste", Price: 50, Tax: 0.1}) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

### 2. Listar Pedidos:
Use a seguinte **query** para listar os pedidos:

```graphql
query ListOrders {
  ListOrders {
    id
    Price
    FinalPrice
  }
}
```

---

## Testes Via gRPC

### Pré-requisitos:
- **Evans**: Caso não tenha instalado o Evans (cliente para gRPC), siga as instruções de instalação: [Evans GitHub](https://github.com/ktr0731/evans).

### Como usar:
1. Acesse o pacote `pb` via Evans.
2. Navegue até o serviço `OrderService`.
3. Selecione a ação desejada (criar ou listar pedidos).