# Projeto de Tracing com OpenTelemetry e Zipkin

Este projeto demonstra como implementar tracing distribuído utilizando OpenTelemetry (OTel) para coletar dados de telemetria de aplicativos e Zipkin para visualizar esses dados.

## Como Usar

### Pré-requisitos

- Docker e Docker Compose instalados na sua máquina.

### Iniciando os Serviços com Docker Compose

1. Clone este repositório e navegue até a pasta raiz do projeto.
2. Dentro do service_b edite o arquivo .env com a sua chave de api para o site weather api
   * https://www.weatherapi.com/login.aspx
3. Execute o seguinte comando para iniciar todos os serviços, incluindo o Zipkin e o OpenTelemetry Collector:

```
docker-compose up --build
```

Isso irá construir as imagens dos serviços (se necessário) e iniciar os contêineres.

### Fazendo uma Requisição POST

Para enviar um CEP e receber informações sobre o clima correspondente à localidade desse CEP, use a seguinte requisição POST:

- **URL**: `http://localhost:8080/validateCEP`
- **Método**: POST
- **Corpo da Requisição como exemplo**:

```
{ "cep": "29902555" }
```

### Visualizando os Resultados de Tracing

Para ver os traces gerados pelas suas requisições:

1. Acesse o Zipkin em seu navegador: `http://localhost:9411/zipkin/`.
2. Na interface do Zipkin, você pode buscar por traces, visualizar detalhes dos spans, tempos de resposta e a estrutura do trace.

Os traces irão mostrar a jornada da sua requisição desde o Serviço A (validação do CEP) até o Serviço B (obtenção de informações do clima), incluindo chamadas individuais para APIs externas, permitindo uma análise detalhada do fluxo e desempenho da sua requisição.

***service_a, service_b*** servem como filtro caso precise filtrar os logs

- totaloperation - mostra o tempo total da chamada
- api_cep_call - mostra o tempo da chamada para api de CEP
- api_weather_call - mostra o tempo da chamada para a api de clima




