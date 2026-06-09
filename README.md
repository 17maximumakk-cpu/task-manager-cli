# Task Manager CLI

Утилита командной строки для управления задачами с поддержкой приоритетов, сроков, фильтрации и экспорта в JSON/CSV.

## Возможности
- Добавление, выполнение, удаление задач
- Приоритеты: `high`, `medium`, `low`
- Срок выполнения (due date)
- Цветной вывод в терминале
- Сохранение задач в `tasks.json`
- Экспорт в CSV

## Установка и запуск

### Python
```bash
cd python
python task_cli.py add "Купить молоко" --priority high --due 2025-12-31
python task_cli.py list
python task_cli.py done 1
python task_cli.py export mytasks.csv
