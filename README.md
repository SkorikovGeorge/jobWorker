# Job Worker System

## Описание

Этот проект представляет собой легковесную, отказоустойчивую систему для обработки фоновых задач (job worker).  Система реализована на Go и использует Redis в качестве очереди сообщений и хранилища статусов задач.  Поддерживается параллельная обработка задач несколькими воркерами, ретраи (повторные попытки выполнения) задач при ошибках, таймауты на выполнение задач, а также возможность приостановки/возобновления обработки задач.  Реализована приоритетная очередь задач на базе Redis Sorted Sets.  Развертывание осуществляется с помощью Docker Compose.

**Основные компоненты:**

*   **API сервер:**  Принимает HTTP запросы на добавление новых задач в очередь и получение статуса задач.
*   **Воркеры:**  Извлекают задачи из очереди Redis и выполняют их.
*   **Redis:**  Используется как очередь сообщений (Redis Sorted Sets) и для хранения статусов задач (Redis Hashes).

## API

API сервер предоставляет следующие endpoints:

*   **`POST /jobs`** - Добавить новую задачу в очередь.

    **Request Body (JSON):**

    ```json
    {
        "priority": 5
    }
    ```
    *   `priority` (int, *optional*): Приоритет задачи (от 1 до 10, где 1 - наивысший приоритет). Если не указан, используется значение по умолчанию (5).

    **Response (JSON):**

    ```json
    {
        "id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
        "status": "pending"
    }
    ```

    *   `id`:  Уникальный идентификатор задачи (UUID).
    *   `status`:  Начальный статус задачи ("pending").


*   **`GET /jobs/{job_id}`** - Получить статус задачи по ID.

    **Path Parameters:**

    *   `job_id`: Уникальный идентификатор задачи (UUID).

    **Response (JSON):**

    ```json
    {
      "status": "<status>"
    }
    ```

    *  `<status>`: текущий статус


*   **`GET /jobs/pause`** - Приостановить обработку задач.

    **Response (JSON):**

    ```json
    {
        "message": "Job processing paused"
    }
    ```

*   **`GET /jobs/resume`** - Возобновить обработку задач.


    **Response (JSON):**

    ```json
    {
        "message": "Job processing resumed"
    }
    ```


## Запуск в Docker Compose

**Запустите систему с помощью Docker Compose:**

    ```bash
    docker compose up --build
    ```

**API сервер по умолчанию будет доступен по адресу `localhost:8080`**
