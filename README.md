# Пример Backend-Сервиса на Golang

Этот проект представляет собой простой backend-сервис, написанный на языке программирования Golang. Сервис открывает два порта: 8080 для HTTP и 50054 для gRPC. Он реализует graceful shutdown для безопасного завершения работы.

## Запуск

Для запуска сервиса, вам потребуется установленный Go.

1. Клонируйте репозиторий:

```bash
git clone https://github.com/Sveboo/server-ads.git
```

2. Перейдите в директорию проекта:

```bash
cd server-ads
```

3. Соберите и запустите сервис:

```bash
go build -o backend .
./backend
```

После выполнения этих шагов, ваш backend-сервис будет доступен на портах 8080 и 50054.

## Зависимости

Для установки зависимостей проекта, используйте `go mod`:

```bash
go mod download
```


## gRPC API

Сервис также предоставляет gRPC API с использованием протокола Protocol Buffers. Для работы с gRPC API необходимо сгенерировать клиентский код на соответствующем языке.

## Graceful Shutdown

Сервис использует graceful shutdown для безопасного завершения работы. Это позволяет завершить текущие запросы и корректно высвободить ресурсы перед остановкой сервера.

## Лицензия

Этот проект лицензируется в соответствии с лицензией MIT.