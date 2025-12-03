# Трассировка и мониторинг

Проект интегрирован с OpenTelemetry для трассировки и мониторинга, что позволяет собирать информацию для отладки и анализа производительности.

## Компоненты

### 1. OpenTelemetry
- **Трассировка**: Отслеживание запросов через систему
- **Метрики**: Сбор метрик производительности
- **Логи**: Интеграция с логами (опционально)

### 2. Jaeger
- **Визуализация трассировок**: UI для просмотра распределенных трассировок
- **Анализ производительности**: Поиск узких мест в системе
- **Отладка**: Просмотр полного пути запроса

### 3. Prometheus
- **Сбор метрик**: Автоматический сбор метрик из приложения
- **Хранение**: Временное хранение метрик
- **Запросы**: PromQL для анализа метрик

## Запуск

### 1. Запуск всех сервисов

```bash
docker-compose up -d
```

Это запустит:
- Backend с трассировкой
- Jaeger UI (порт 16686)
- Prometheus (порт 9091)
- PostgreSQL

### 2. Доступ к интерфейсам

- **Jaeger UI**: http://localhost:16686
- **Prometheus UI**: http://localhost:9091
- **Backend API**: http://localhost:8000
- **Metrics endpoint**: http://localhost:8000/metrics

## Использование

### Просмотр трассировок в Jaeger

1. Откройте http://localhost:16686
2. Выберите сервис `tutoring-platform`
3. Нажмите "Find Traces"
4. Выберите trace для детального просмотра

### Просмотр метрик в Prometheus

1. Откройте http://localhost:9091
2. Перейдите в раздел "Graph"
3. Введите PromQL запрос, например:
   ```
   http_requests_total
   http_request_duration_seconds
   ```

### Доступные метрики

#### HTTP метрики:
- `http_requests_total` - общее количество HTTP запросов
- `http_request_duration_seconds` - длительность HTTP запросов
- `http_request_size_bytes` - размер входящих запросов
- `http_response_size_bytes` - размер исходящих ответов

Все метрики имеют labels:
- `http.method` - HTTP метод (GET, POST, etc.)
- `http.route` - путь запроса
- `http.status_code` - HTTP статус код

## Конфигурация

### Переменные окружения

```bash
# Имя сервиса для трассировки
SERVICE_NAME=tutoring-platform

# Версия сервиса
SERVICE_VERSION=1.0.0

# Endpoint Jaeger collector
JAEGER_ENDPOINT=http://jaeger:14268/api/traces
```

### Настройка Prometheus

Конфигурация находится в `tools/prometheus/prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'tutoring-platform'
    static_configs:
      - targets: ['backend:9090']
```

## Использование в коде

### Создание span для трассировки

```go
import (
    "data_base_project/tracing"
    "go.opentelemetry.io/otel/trace"
)

func myFunction(ctx context.Context) {
    tracer := tracing.GetTracer("my-component")
    ctx, span := tracer.Start(ctx, "my-operation")
    defer span.End()
    
    // Ваш код здесь
    span.SetAttributes(
        attribute.String("key", "value"),
    )
}
```

### Добавление метрик

```go
import (
    "data_base_project/tracing"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/metric"
)

func recordMetric() {
    meter := tracing.GetMeter("my-component")
    
    counter, _ := meter.Int64Counter(
        "my_counter",
        metric.WithDescription("My counter"),
    )
    
    counter.Add(context.Background(), 1,
        metric.WithAttributes(
            attribute.String("label", "value"),
        ),
    )
}
```

## Отладка

### Проверка работы трассировки

1. Сделайте несколько запросов к API
2. Проверьте Jaeger UI - должны появиться traces
3. Если traces не появляются:
   - Проверьте логи backend: `docker-compose logs backend`
   - Убедитесь, что Jaeger доступен: `docker-compose ps jaeger`
   - Проверьте переменные окружения

### Проверка метрик

1. Откройте http://localhost:8000/metrics
2. Должны быть видны метрики OpenTelemetry
3. В Prometheus проверьте targets: http://localhost:9091/targets

### Типичные проблемы

#### Traces не появляются в Jaeger
- Проверьте, что `JAEGER_ENDPOINT` правильный
- Убедитесь, что Jaeger контейнер запущен
- Проверьте сеть Docker: `docker network inspect <network>`

#### Метрики не собираются
- Проверьте endpoint `/metrics`
- Убедитесь, что Prometheus может достучаться до backend
- Проверьте конфигурацию Prometheus

## Production рекомендации

1. **Sampling**: Измените sampler с `AlwaysSample` на более умный (например, `TraceIDRatioBased`)
2. **Retention**: Настройте retention для Jaeger и Prometheus
3. **Security**: Добавьте аутентификацию для Jaeger и Prometheus UI
4. **Resource limits**: Установите лимиты ресурсов для контейнеров
5. **Backup**: Настройте бэкапы для данных трассировок

## Дополнительные ресурсы

- [OpenTelemetry Go Documentation](https://opentelemetry.io/docs/instrumentation/go/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [Prometheus Documentation](https://prometheus.io/docs/)

