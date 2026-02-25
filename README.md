# AWS Stresser Observability

## 📋 Escopo Geral do Projeto

Este projeto é uma **plataforma completa de observabilidade** que demonstra a implementação de um stack moderno de monitoramento, logging e tracing distribuído. O sistema consiste em uma aplicação Go de stress test de CPU integrada com ferramentas de observabilidade, utilizando LocalStack para simular serviços AWS S3 para armazenamento de logs e traces.

### Objetivo Principal

Demonstrar a implementação prática dos três pilares da observabilidade:
- **Métricas**: Coleta e visualização de métricas de performance via Prometheus
- **Logs**: Agregação e consulta de logs estruturados via Loki
- **Traces**: Rastreamento distribuído de requisições via Tempo e OpenTelemetry

## 🏗️ Arquitetura

```
┌─────────────────┐
│  Stresser App   │ ──► Gera métricas, logs e traces
│   (Go + OTEL)   │
└────────┬────────┘
         │
         ├──► Métricas ──► Prometheus ──► Grafana
         │
         ├──► Logs ──► Promtail ──► Loki ──► S3 (LocalStack)
         │
         └──► Traces ──► OTel Collector ──► Tempo ──► S3 (LocalStack)
```

Para visualizar a arquitetura do projeto entrar em: `docs/aws-stresser-observability-architecture.png`

## 🛠️ Ferramentas Utilizadas

### Aplicação Principal
- **Go**: Linguagem de programação para a aplicação de stress test
- **OpenTelemetry**: Instrumentação para tracing distribuído
- **Prometheus Client**: Exportação de métricas customizadas

### Stack de Observabilidade

#### Métricas
- **Prometheus**: Sistema de monitoramento e time-series database
  - Scraping de métricas a cada 5 segundos
  - Armazenamento local em volume persistente
  - Regras de alertas configuradas

#### Logs
- **Loki**: Sistema de agregação de logs inspirado no Prometheus
  - Armazenamento no S3 (LocalStack) via bucket `loki-chunks`
  - Schema v11 com BoltDB Shipper
  - Retenção configurável de logs
  
- **Promtail**: Agente de coleta de logs
  - Coleta logs dos containers Docker
  - Parsing e labeling automático
  - Push para Loki via HTTP

#### Traces
- **Tempo**: Backend de tracing distribuído da Grafana
  - Armazenamento no S3 (LocalStack) via bucket `tempo-traces`
  - Recepção via protocolo OTLP
  - Retenção de 1 hora para traces
  - Integração com Loki e Prometheus

- **OpenTelemetry Collector**: Coletor e processador de telemetria
  - Recebe traces via OTLP
  - Processamento em batch
  - Export para Tempo

#### Visualização
- **Grafana**: Plataforma de visualização e analytics
  - Json Provisioning dashboards
  - Datasources pré-configurados (Prometheus, Loki, Tempo)
  - Correlação automática entre traces, logs e métricas

#### Alertas
- **Alertmanager**: Gerenciamento de alertas do Prometheus
  - Roteamento de alertas para email
  - Agrupamento e deduplicação
  - Integração com MailHog para testes

- **MailHog**: Servidor SMTP de teste
  - Captura emails sem envio real
  - Interface web para visualização

#### Infraestrutura
- **LocalStack**: Emulador de serviços AWS
  - Simula S3 para armazenamento de logs e traces
  - Buckets: `loki-chunks`, `tempo-traces`, `stresser-logs`
  - Persistência de dados em volume
  - Endpoint: `http://localstack:4566`

- **Docker Compose**: Orquestração de containers
  - 11 serviços integrados
  - Volumes persistentes
  - Health checks e dependências

## 📊 Dashboard Grafana - Provisionamento via JSON

### Estrutura de Provisionamento

O Grafana utiliza o sistema de **provisioning automático** para configurar datasources e dashboards na inicialização:

```
grafana/provisioning/
├── datasources/
│   └── datasources.yml          # Configuração dos datasources
└── dashboards/
    ├── dashboard-config.yml     # Configuração do provider
    └── stress-app-and-logs.json # Dashboard completo
```

### Datasources Provisionados

**1. Prometheus** (Default)
- URL: `http://prometheus:9090`
- Tipo: Métricas time-series
- Uso: Visualização de métricas de CPU, requests, workers

**2. Loki**
- URL: `http://loki:3100`
- Tipo: Logs agregados
- Configuração: Max 1000 linhas por query

**3. Tempo**
- URL: `http://tempo:3200`
- Tipo: Traces distribuídos
- Integrações especiais:
  - **Traces → Logs**: Correlação automática com Loki usando tags de serviço
  - **Traces → Metrics**: Queries automáticas no Prometheus
  - **Service Map**: Visualização de dependências entre serviços
  - **Node Graph**: Grafo de spans e latências

## 💾 Armazenamento no S3 Simulado (LocalStack)

O LocalStack simula o serviço S3 da AWS localmente, permitindo desenvolvimento e testes sem custos

### Inicialização dos Buckets

O serviço `s3-init` cria automaticamente os buckets necessários:

```bash
aws s3 mb s3://loki-chunks
aws s3 mb s3://tempo-traces
aws s3 mb s3://stresser-logs
```

### Integração Loki → S3

**Fluxo de Dados:**
1. Promtail coleta logs dos containers Docker
2. Promtail envia logs para Loki via HTTP
3. Loki processa e indexa os logs
4. Loki armazena chunks no bucket `loki-chunks` do LocalStack
5. Schema BoltDB Shipper gerencia índices localmente
6. Dados persistem no volume `localstack-data`

### Integração Tempo → S3

**Fluxo de Dados:**
1. Aplicação Go gera spans OpenTelemetry
2. Spans são enviados para OTel Collector via gRPC (porta 4317)
3. OTel Collector processa em batch e exporta para Tempo
4. Tempo recebe traces via OTLP gRPC (porta 4317)
5. Tempo armazena traces no bucket `tempo-traces` do LocalStack
6. WAL (Write-Ahead Log) temporário em `/tmp/tempo/wal`
7. Compactação automática com retenção de 1 hora

### Benefícios do S3 Simulado

1. **Desenvolvimento Local**: Sem necessidade de conta AWS
2. **Custo Zero**: Sem cobranças por armazenamento ou requisições
3. **Velocidade**: Latência mínima em ambiente local
4. **Persistência**: Dados mantidos entre reinicializações
5. **Realismo**: API compatível com S3 real
6. **Migração Fácil**: Trocar endpoint para S3 real em produção

## 🚀 Como Executar

### Pré-requisitos
- Docker e Docker Compose instalados
- Portas disponíveis: 3000, 8080, 9090, 3100, 3200, 4566, 9093, 8025

### Inicialização

```bash
# Subir todos os serviços
docker-compose up --build

# Verificar status
docker-compose ps
```
### Acessos

- **Stresser App**: http://localhost:8080
- **Grafana**: http://localhost:3000
- **Prometheus**: http://localhost:9090
- **Alertmanager**: http://localhost:9093
- **MailHog**: http://localhost:8025
- **LocalStack**: http://localhost:4566

## 📈 Métricas Disponíveis

A aplicação exporta as seguintes métricas Prometheus:

- `stress_level`: Nível atual de stress (0-100)
- `stress_cpu_workers`: Número de goroutines ativas
- `stress_changes_total`: Total de mudanças de nível
- `http_requests_total`: Total de requisições HTTP
- `http_request_duration_seconds`: Latência de requisições
- `estimated_cost_usd`: Custo simulado em USD

## 🔔 Sistema de Alertas

**Regra Configurada:**
- **Nome**: stress-acima-do-limite
- **Condição**: `stress_level > 80`
- **Duração**: 10 segundos
- **Severidade**: Critical
- **Ação**: Email via MailHog

## 📦 Volumes Persistentes

- `grafana-data`: Dashboards e configurações do Grafana
- `prometheus-data`: Time-series database do Prometheus
- `localstack-data`: Buckets S3 simulados (logs e traces)

## 🔧 Configurações Importantes

### Limites de Recursos
```yaml
stresser-app:
  deploy:
    resources:
      limits:
        cpus: "2.0"
        memory: "512M"
```

### Retenção de Dados
- **Tempo**: 1 hora de traces
- **Loki**: Configurável via schema
- **Prometheus**: Padrão (15 dias)

## 🎯 Casos de Uso

1. **Aprendizado**: Entender observabilidade na prática
2. **Testes**: Simular carga e observar comportamento
3. **Desenvolvimento**: Prototipar integrações com S3

## 🌟 Destaques Técnicos

- ✅ Observabilidade completa (métricas, logs, traces)
- ✅ Provisionamento automático de dashboards
- ✅ Armazenamento em S3 simulado
- ✅ Correlação automática entre sinais
- ✅ Sistema de alertas funcional
- ✅ Interface web interativa
- ✅ Instrumentação OpenTelemetry
- ✅ Logs estruturados em JSON
- ✅ Persistência de dados
