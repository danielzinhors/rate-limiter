# O que é isso?

Este é um middleware limitador de taxa configurável para Golang que bloqueia requisições após o limite estipulado de tentaivas tanto por IP com por Token.

# Como funciona?

O RateLimiter utiliza o redis para armazenar os acessos dos IP e Tokens, a fim fr bloquea-los em caso de ultrapassar o limite de rqueisições permitidas por segundo, 
deixando-os bloqueados pelo tempo determinado nas configurações.

# Como testar?

Execute-o com o docker compose:

No bash execute o comando
```bash
docker compose up -d
```

Ative p modo de depuração para poder acompanhar os logs com:
No bash execute o comando
```bash
docker compose logs server --follow 
```

Você pode enviar solicitações de qualquer browser ou outro aplicativo que realize requisições HTTP (REST Client como Postman, insominia etc)
utilizando o endpoint `http://localhost:8080` com uma chamada GET ou usar um testador de estresse, como o postman runner. 
Ou utilizando os arquivos .http na pasta api do projeto

# Como configurar o projeto?

## Variáveis ​​ambientais

Você pode configurar as variáveis de ambiente apenas editando o arquivo `.env`:

|Value|Type|Description|Default Value|
|---|---|---|---|
|MAX_REQUESTS_RATE_LIMITER_IP|integer|Solicitações por segundo permitidas para um IP.|100|

|BLOCK_TIME_RATE_LIMITER_IP|integer|Tempo de bloqueio em milissegundos para IPs que atingem sua cota de solicitações.|1000|

|MAX_REQUESTS_RATE_LIMITER_TOKEN|integer|Solicitações por segundo permitidas para um token. Isto tem prioridade sobre a configuração IP.|200|

|BLOCK_TIME_RATE_LIMITER_TOKEN|integer|Tempo de bloqueio em milissegundos para tokens que atinjam sua cota de solicitação. Isto tem prioridade sobre a configuração IP.|500|

|MAX_REQUESTS_RATE_LIMITER_TOKEN_ABC|integer|Solicitações por segundo permitidas para o token "ABC". Se não for definido, usará MAX_REQUESTS_RATE_LIMITER_TOKEN. |-|

|BLOCK_TIME_RATE_LIMITER_TOKEN_ABC|integer|Tempo de bloqueio em milissegundos para o token "ABC". Se não for definido, usará BLOCK_TIME_RATE_LIMITER_TOKEN. |-|

|DEBUG_RATE_LIMITER|boolean|Executa em modo de depuração e mensagens são exibidas bash.|false|

|USE_RATE_LIMITER_REDIS|boolean|Usa o Adpter de Storage do Redis.|false|

|ADDRESS_RATE_LIMITER_REDIS|string|Endereço para o Adpter de Storage do Redis.|-|

|PASSWORD_RATE_LIMITER_REDIS|string|Password para o Adpter de Storage do Redis.|-|

|DB_RATE_LIMITER_REDIS|integer|Database para o Adpter de Storage do Redis.|-|
