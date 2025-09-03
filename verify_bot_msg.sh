#!/bin/bash

echo "=== 验证bot_msg任务执行 ==="

# 启动服务
echo "1. 启动服务..."
go run main.go -mode=dev &
SERVICE_PID=$!

sleep 3

echo "2. 测试立即执行任务..."
curl -X POST http://localhost:8080/api/task/test \
  -H "Content-Type: application/json" \
  -d '{
    "task_type": "bot_msg",
    "payload": {"msg_type": "test", "content": "测试bot_msg任务"}
  }'

echo "3. 测试cron定时任务..."
curl -X POST http://localhost:8080/api/task/create \
  -H "Content-Type: application/json" \
  -d '{
    "task_type": "bot_msg",
    "cron_expression": "*/2 * * * *",
    "payload": {"msg_type": "cron_test", "content": "每2分钟测试"}
  }'

echo "4. 等待30秒后检查日志..."
sleep 30

echo "5. 检查调度器条目..."
# 可以通过API或直接查看Redis

echo "=== 验证完成 ==="
kill $SERVICE_PID 2>/dev/null || true