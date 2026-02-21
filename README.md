// ...existing code...
# green_daemon

green_daemon is a lightweight system process monitoring daemon written in Go that exposes Prometheus metrics for visualization in Grafana. A React (Vite) frontend embeds the Grafana panel for simple UI display.

---

## 🧱 Tech stack

- Go — Metrics daemon
- Prometheus — Metrics collection & storage
- Grafana — Visualization & dashboards
- React (Vite) — Frontend embedding Grafana panel

---

## 🏗 Architecture

green_daemon → Prometheus → Grafana → React frontend

1. The Go daemon polls running processes.
2. Metrics are exposed at `/metrics`.
3. Prometheus scrapes metrics periodically.
4. Grafana visualizes metrics.
5. React embeds the Grafana panel for UI.

---

## 🚀 Run everything locally

Open 4 separate terminals.

### 1) Start backend (green_daemon)

```bash
cd backend
go mod tidy
go run .
```

Metrics endpoint: http://localhost:8080/metrics

---

### 2) Start Prometheus

From project root:

```bash
prometheus --config.file=prometheus/prometheus.yml
```

Prometheus UI: http://localhost:9090  
Verify targets: http://localhost:9090/targets  
Ensure the backend target is UP.

---

### 3) Start Grafana

If installed with Homebrew:

```bash
brew services start grafana
```

Grafana UI: http://localhost:3000  
Default login: admin / admin

Ensure Prometheus datasource is configured as: http://localhost:9090

To embed Grafana in an iframe, ensure allow_embedding = true is enabled in Grafana config.

---

### 4) Start frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend: http://localhost:5173

---

## 🧪 Verify

- http://localhost:8080/metrics → Metrics visible
- http://localhost:9090/targets → Backend target UP
- Grafana dashboard displays data
- React frontend displays embedded Grafana panel

---

## 📊 Example PromQL queries

Top 10 memory consuming processes:

```promql
topk(10, mem_utlization)
```

Top 5 CPU time consuming processes:

```promql
topk(5, cpu_time_total)
```

All memory metrics:

```promql
mem_utlization
```

---

## 🛑 Stop services

- Stop backend / Prometheus: press CTRL+C in their terminals
- Stop Grafana (Homebrew): `brew services stop grafana`
- Stop frontend: press CTRL+C

---

## ⚠️ Notes

- Metrics are labeled by process name and PID, which may cause high cardinality.
- For production, consider removing PID labeling or limiting to top N processes.

---

## 📌 Future improvements

- Reduce metric cardinality
- Add system-wide CPU & memory metrics
- Add alerting rules
- Dockerize services
- Kubernetes deployment
- Add CI/CD workflow

---
