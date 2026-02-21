export default function App() {
  return (
    <div style={{ fontFamily: "system-ui", padding: 24 }}>
      <h1>GreenOps Monitoring</h1>

      <div
        style={{
          marginTop: 32,
          border: "1px solid #222",
          borderRadius: 12,
          overflow: "hidden",
        }}
      >
        <iframe
          src="http://127.0.0.1:3000/d-solo/adxsxbs/requests-per-second?orgId=1&panelId=panel-1&from=now-5m&to=now&refresh=5s"
          width="100%"
          height="400"
          frameBorder="0"
          title="Grafana RPS"
        />
      </div>
    </div>
  );
}