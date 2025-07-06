# 🕒 Cron-Job-Scheduler (Go + Redis)

A lightweight, modular cron job scheduler built in Go, using Redis for persistent scheduling. Accepts HTTP-based task submissions and allows developers to define custom task logic easily. Currently, it performs a basic **ping check** to verify if a submitted website is online — but you can **extend this system with your own task logic**.

---

## 📌 Features

* ⚡ Accepts HTTP POST requests to create jobs dynamically
* 🧠 Redis-backed storage for task metadata and sorted scheduling
* 👷‍♂️ Worker pool to handle concurrent task execution
* ⏱️ Delayed execution using Redis Sorted Sets (ZSET)
* 🔄 Task queue that re-fetches and evaluates task timing continuously
* 🔧 **Modular architecture** — add your own task logic (currently a simple status check)
* 🧹 Graceful cancellation and rescheduling of jobs

---

## 🛠 Modular Task Execution

This project is designed to be modular. You can easily plug in **custom task behaviors** by modifying the `get_data` function or extending the worker logic.

### ✅ Current Behavior:

* It sends an HTTP GET to the submitted URL and checks the status code.
* If the site responds (i.e., `200 OK`), the job is considered successful.

### 🧩 You Can Extend It To:

* Trigger webhooks
* Write logs to a file or external service
* Push metrics to Prometheus
* Retry failed requests
* Send emails or SMS alerts
* Run shell commands or Docker containers

> All tasks follow a pattern that can be replaced or extended without breaking the scheduling system.

---

## 📥 API Usage

### 🔧 POST `/post`

Create a new scheduled task.

#### Body Parameters (JSON):

```json
{
  "url": "http://example.com/task",
  "delay": 30
}
```

| Field | Type    | Description                       |
| ----- | ------- | --------------------------------- |
| url   | string  | The URL to hit when the task runs |
| delay | integer | Delay in seconds before execution |

#### Example Curl:

```bash
curl -X POST http://localhost:5000/post \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8000/test", "delay": 15}'
```

---

## 🧠 How It Works

1. Client sends a job with a delay and target URL.
2. Task is stored in Redis (hash + sorted set).
3. Worker pool checks when the task is due.
4. Executes `get_data(url)` — currently, a ping/status check.
5. Logs the result and listens for more tasks.

---

## ⚙️ Project Structure

```
.
├── main.go        # Handles HTTP, Redis client, job submission
├── worker.go      # Worker pool and scheduling logic
├── redis.go       # Redis task handling (HSET, ZADD, etc.)
├── utils.go       # Utility functions like ID generation
```

---

## 📦 Requirements

* Go 1.19+
* Redis 6+

---

## ▶️ Running the Project

1. Start Redis locally:

   ```bash
   redis-server
   ```

2. Run the server:

   ```bash
   go run main.go redis.go worker.go utils.go
   ```

3. Submit a task using curl or Postman.

---

## 📈 Example Output

```bash
workerid : 2 success: 200
workerid : 1 success: 200
```

---

## ✅ To Do

* Retry mechanism for failed tasks
* Job deletion
* Add repetitive tasks

---

## 📄 License

MIT License. Fork it, change it, deploy it.

---

## 🙌 Author

Made by [@Omkardalvi01](https://github.com/Omkardalvi01)
