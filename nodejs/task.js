#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const TASKS_FILE = path.join(__dirname, 'tasks.json');

const colors = {
    red: '\x1b[91m',
    green: '\x1b[92m',
    yellow: '\x1b[93m',
    blue: '\x1b[94m',
    reset: '\x1b[0m'
};

function colorize(text, color) {
    return `${colors[color] || ''}${text}${colors.reset}`;
}

function loadTasks() {
    if (!fs.existsSync(TASKS_FILE)) return [];
    const data = fs.readFileSync(TASKS_FILE, 'utf8');
    return JSON.parse(data);
}

function saveTasks(tasks) {
    fs.writeFileSync(TASKS_FILE, JSON.stringify(tasks, null, 2), 'utf8');
}

function addTask(description, priority = 'medium', dueDate = '') {
    const tasks = loadTasks();
    const newId = tasks.length > 0 ? Math.max(...tasks.map(t => t.id)) + 1 : 1;
    tasks.push({
        id: newId,
        description,
        status: 'pending',
        priority,
        due_date: dueDate,
        created_at: new Date().toISOString()
    });
    saveTasks(tasks);
    console.log(colorize(`✓ Задача ${newId} добавлена.`, 'green'));
}

function listTasks(statusFilter = null) {
    let tasks = loadTasks();
    if (statusFilter) tasks = tasks.filter(t => t.status === statusFilter);
    if (tasks.length === 0) {
        console.log(colorize('Нет задач.', 'yellow'));
        return;
    }
    tasks.forEach(t => {
        const statusIcon = t.status === 'done' ? '✓' : '◻';
        const prioColor = { high: 'red', medium: 'yellow', low: 'green' }[t.priority] || 'reset';
        const due = t.due_date ? ` [до ${t.due_date}]` : '';
        console.log(`${statusIcon} [${t.id}] ${colorize(t.description, prioColor)}${due}`);
    });
}

function completeTask(taskId) {
    const tasks = loadTasks();
    const task = tasks.find(t => t.id === taskId);
    if (!task) {
        console.log(colorize(`Ошибка: задача ${taskId} не найдена.`, 'red'));
        return;
    }
    task.status = 'done';
    saveTasks(tasks);
    console.log(colorize(`✓ Задача ${taskId} выполнена.`, 'green'));
}

function deleteTask(taskId) {
    let tasks = loadTasks();
    const newTasks = tasks.filter(t => t.id !== taskId);
    if (newTasks.length === tasks.length) {
        console.log(colorize(`Ошибка: задача ${taskId} не найдена.`, 'red'));
    } else {
        saveTasks(newTasks);
        console.log(colorize(`✓ Задача ${taskId} удалена.`, 'green'));
    }
}

function exportCsv(filename) {
    const tasks = loadTasks();
    const headers = ['id', 'description', 'status', 'priority', 'due_date', 'created_at'];
    const rows = tasks.map(t => headers.map(h => t[h] || '').join(','));
    const csv = [headers.join(','), ...rows].join('\n');
    fs.writeFileSync(filename, csv, 'utf8');
    console.log(colorize(`Экспортировано ${tasks.length} задач в ${filename}`, 'blue'));
}

// --- CLI аргументы ---
const args = process.argv.slice(2);
const command = args[0];
switch (command) {
    case 'add':
        const desc = args[1];
        if (!desc) { console.log('Использование: task.js add "описание" [--priority high|low] [--due YYYY-MM-DD]'); process.exit(1); }
        let priority = 'medium';
        let due = '';
        for (let i = 2; i < args.length; i++) {
            if (args[i] === '--priority' && args[i+1]) priority = args[++i];
            if (args[i] === '--due' && args[i+1]) due = args[++i];
        }
        addTask(desc, priority, due);
        break;
    case 'list':
        const status = args[1] === '--status' ? args[2] : null;
        listTasks(status);
        break;
    case 'done':
        const idDone = parseInt(args[1]);
        if (isNaN(idDone)) { console.log('Укажите ID задачи'); process.exit(1); }
        completeTask(idDone);
        break;
    case 'rm':
        const idRm = parseInt(args[1]);
        if (isNaN(idRm)) { console.log('Укажите ID задачи'); process.exit(1); }
        deleteTask(idRm);
        break;
    case 'export':
        const filename = args[1] || 'tasks.csv';
        exportCsv(filename);
        break;
    default:
        console.log(`
Доступные команды:
  add "описание" [--priority high|medium|low] [--due YYYY-MM-DD]
  list [--status pending|done]
  done <id>
  rm <id>
  export [filename.csv]
        `);
}
