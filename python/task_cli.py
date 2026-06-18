#!/usr/bin/env python3
"""
Task Manager CLI – управление задачами в терминале
Поддерживает: JSON/CSV, приоритеты, сроки, цветной вывод
"""

import json
import csv
import os
import sys
import argparse
from datetime import datetime
from typing import List, Dict, Any

TASKS_FILE = "tasks.json"
COLORS = {
    "red": "\033[91m",
    "green": "\033[92m",
    "yellow": "\033[93m",
    "blue": "\033[94m",
    "reset": "\033[0m"
}

def colorize(text: str, color: str) -> str:
    return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"

def load_tasks() -> List[Dict[str, Any]]:
    if not os.path.exists(TASKS_FILE):
        return []
    with open(TASKS_FILE, "r", encoding="utf-8") as f:
        return json.load(f)

def save_tasks(tasks: List[Dict[str, Any]]):
    with open(TASKS_FILE, "w", encoding="utf-8") as f:
        json.dump(tasks, f, indent=2, ensure_ascii=False)

def add_task(description: str, priority: str = "medium", due_date: str = ""):
    tasks = load_tasks()
    new_id = max([t["id"] for t in tasks], default=0) + 1
    tasks.append({
        "id": new_id,
        "description": description,
        "status": "pending",
        "priority": priority,
        "due_date": due_date,
        "created_at": datetime.now().isoformat()
    })
    save_tasks(tasks)
    print(colorize(f"✓ Задача {new_id} добавлена.", "green"))

def list_tasks(status_filter: str = None):
    tasks = load_tasks()
    if status_filter:
        tasks = [t for t in tasks if t["status"] == status_filter]
    if not tasks:
        print(colorize("Нет задач.", "yellow"))
        return
    for t in tasks:
        status_icon = "✓" if t["status"] == "done" else "◻"
        prio_color = {"high": "red", "medium": "yellow", "low": "green"}.get(t["priority"], "reset")
        due = f" [до {t['due_date']}]" if t["due_date"] else ""
        print(f"{status_icon} [{t['id']}] {colorize(t['description'], prio_color)}{due}")

def complete_task(task_id: int):
    tasks = load_tasks()
    for t in tasks:
        if t["id"] == task_id:
            t["status"] = "done"
            save_tasks(tasks)
            print(colorize(f"✓ Задача {task_id} выполнена.", "green"))
            return
    print(colorize(f"Ошибка: задача {task_id} не найдена.", "red"))

def delete_task(task_id: int):
    tasks = load_tasks()
    new_tasks = [t for t in tasks if t["id"] != task_id]
    if len(new_tasks) == len(tasks):
        print(colorize(f"Ошибка: задача {task_id} не найдена.", "red"))
    else:
        save_tasks(new_tasks)
        print(colorize(f"✓ Задача {task_id} удалена.", "green"))

def export_csv(filename: str):
    tasks = load_tasks()
    with open(filename, "w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=["id", "description", "status", "priority", "due_date", "created_at"])
        writer.writeheader()
        writer.writerows(tasks)
    print(colorize(f"Экспортировано {len(tasks)} задач в {filename}", "blue"))

def main():
    parser = argparse.ArgumentParser(description="Менеджер задач CLI")
    subparsers = parser.add_subparsers(dest="command", required=True)

    add_p = subparsers.add_parser("add", help="Добавить задачу")
    add_p.add_argument("description", help="Описание задачи")
    add_p.add_argument("--priority", choices=["low", "medium", "high"], default="medium")
    add_p.add_argument("--due", default="", help="Дата в формате ГГГГ-ММ-ДД")

    list_p = subparsers.add_parser("list", help="Показать задачи")
    list_p.add_argument("--status", choices=["pending", "done"], help="Фильтр по статусу")

    done_p = subparsers.add_parser("done", help="Отметить задачу выполненной")
    done_p.add_argument("id", type=int)

    rm_p = subparsers.add_parser("rm", help="Удалить задачу")
    rm_p.add_argument("id", type=int)

    export_p = subparsers.add_parser("export", help="Экспорт в CSV")
    export_p.add_argument("filename", default="tasks.csv", nargs="?")

    args = parser.parse_args()
    if args.command == "add":
        add_task(args.description, args.priority, args.due)
    elif args.command == "list":
        list_tasks(args.status)
    elif args.command == "done":
        complete_task(args.id)
    elif args.command == "rm":
        delete_task(args.id)
    elif args.command == "export":
        export_csv(args.filename)

if __name__ == "__main__":
    main()
