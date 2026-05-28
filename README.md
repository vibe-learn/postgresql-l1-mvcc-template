        # postgresql — ACID и MVCC: как PG не блокирует чтения

        Homework-шаблон для урока **l1_acid_and_mvcc** (ACID и MVCC: как PG не блокирует чтения) на платформе Vibe Learn.

        ## Что делать

        Дано: testcontainers PG. Реализуй на Go демо MVCC:
1) Создай таблицу + индексы.
2) Запусти 4 параллельных goroutine: 1 пишет (UPDATE одной строки в цикле),
   3 читают snapshot этой строки в долгих транзакциях.
3) Показатели до/после: `pg_stat_user_tables.n_dead_tup`, размер таблицы и индексов
   (pg_relation_size).
4) Покажи, что блокировок не возникает (нет ожиданий в pg_locks).
5) Запусти VACUUM и измерь, сколько dead tuples ушло.
Тесты проверят корректность измерений и absence of read locks.

## Контекст (из transfer-задачи урока)

Поступила жалоба от Ops: одна из таблиц в проде, `metrics_minute_aggregates`, выросла
с 80 ГБ до 280 ГБ за месяц. Запросы по ней замедлились в 3 раза. Количество строк
выросло всего в 1.4 раза. Что характерно:

- таблица — массив агрегатов, на каждую (host, metric_name) каждую минуту делается
  UPSERT (`INSERT ... ON CONFLICT DO UPDATE`);
- autovacuum_vacuum_scale_factor — дефолтный 0.2 (20%);
- индексов на таблице 6 штук;
- pg_stat_user_tables показывает n_dead_tup ≈ 3 × n_live_tup;
- последний autovacuum по этой таблице — 4 дня назад.

## Recap из урока

- **MVCC = главное архитектурное свойство PG**: чтения не блокируют записи, потому что каждая транзакция видит свой snapshot.
- **UPDATE никогда не правит строку «на месте».** Создаётся новая версия (xmin=current_xid), старая помечается удалённой (xmax=current_xid).
- **Цена MVCC — dead tuples и bloat.** VACUUM/autovacuum их чистят, но обычный VACUUM не возвращает место ОС — только помечает свободным.
- **Atomicity и durability — побочный эффект WAL.** Write-ahead log + fsync на COMMIT = «или всё, или ничего» при любом crash.
- **Autovacuum включён по умолчанию, но часто отстаёт на write-heavy таблицах.** Уменьшай `autovacuum_vacuum_scale_factor` для горячих таблиц.

        ## Как работать

        1. Платформа Vibe Learn создаёт копию этого репо в твоём GitHub-аккаунте по клику «Начать домашку» на странице урока (через GitHub `/generate`, codecrafters-pattern).
        2. Склонируй копию локально, реализуй TODO в `main.go`, прогони тесты, запушь.
        3. CI (`.github/workflows/ci.yml`) запускает `go vet` + `go test ./...` на каждый push. Платформа слушает результат через webhook от GitHub Actions и обновляет статус домашки на странице урока.

        ## Локальное окружение

        - Go 1.22+
        - Docker + docker-compose — `docker compose up -d` поднимает single-node PostgreSQL 16 на `localhost:5432` с healthcheck. DSN: `postgres://postgres:postgres@localhost:5432/postgres`. Переопределяется через env `DATABASE_URL`.

        ## Запуск

        ```bash
        # Поднять локальный PostgreSQL
        docker compose up -d

        # Прогнать тесты (интеграционный включается через PG_INTEGRATION=1)
        go test ./...
        PG_INTEGRATION=1 go test ./...

        # Запустить main (печатает marker; замени stub на реализацию)
        go run .
        ```

        ## Заметка автора

        Это baseline-шаблон, сгенерированный платформой. Бизнес-сущность задачи (что конкретно реализовать в `main.go`, какие тесты сделать строгими) расширяется по ходу итераций — параллельно с углублением теории урока.
