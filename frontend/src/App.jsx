import { useState, useRef, useEffect } from "react";

export default function App() {
  const [chatOpen, setChatOpen] = useState(false);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [activePanelTitle, setActivePanelTitle] = useState("");
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const panelSrcs = [
    {
      src: "http://127.0.0.1:3000/d-solo/adxsxbs/requests-per-second?orgId=1&panelId=panel-1&from=now-5m&to=now&refresh=5s",
      title: "Requests per second",
    },
    {
      src: "http://localhost:3000/d-solo/adxsxbs/system-process-monitoring?orgId=1&from=1771702145795&to=1771705745795&timezone=browser&panelId=panel-2&__feature.dashboardSceneSolo=true",
      title: "System Process Monitoring",
    },
    {
      src: "http://localhost:3000/d-solo/packets-timing-1/packet-timing-per-process?orgId=1&from=now-1h&to=now&panelId=2",
      title: "Packet Timing per Process",
    },
  ];

  // Calls your backend (which calls Gemini). Keep the API key ONLY on the backend.
  const askGemini = async ({ query, context }) => {
    setLoading(true);
    try {
      const response = await fetch("http://127.0.0.1:3001/api/ai/query", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query, context }),
      });

      const data = await response.json();
      const assistantMessage = {
        role: "assistant",
        content: data.message || data.response || "Unable to analyze the metric",
      };
      setMessages((prev) => [...prev, assistantMessage]);
    } catch (error) {
      setMessages((prev) => [
        ...prev,
        {
          role: "assistant",
          content:
            "Error connecting to Gemini backend. Make sure your server is running on http://127.0.0.1:3001.",
        },
      ]);
    } finally {
      setLoading(false);
    }
  };

  // Click 🤖 -> open chat + auto-call Gemini with panel context
  const triggerAnalysis = async (panelTitle) => {
    setChatOpen(true);
    setActivePanelTitle(panelTitle);

    const autoQuery = `Analyze the spike for panel "${panelTitle}". Give diagnosis, likely causes, checks, and mitigation.`;
    const context = `Grafana panel: ${panelTitle}. You are analyzing monitoring metrics and performance spikes.`;

    // Show the auto-query as a user message (nice UX)
    setMessages((prev) => [...prev, { role: "user", content: autoQuery }]);

    await askGemini({ query: autoQuery, context });
  };

  const handleSendMessage = async () => {
    if (!input.trim() || loading) return;

    const userText = input;
    setMessages((prev) => [...prev, { role: "user", content: userText }]);
    setInput("");

    const context = activePanelTitle
      ? `Grafana panel: ${activePanelTitle}. You are analyzing monitoring metrics and performance spikes.`
      : "Analyzing monitoring metrics and performance spikes.";

    await askGemini({ query: userText, context });
  };

  return (
    <div
      style={{
        fontFamily: "system-ui",
        width: "100vw",
        height: "100vh",
        padding: 0,
        margin: 0,
        overflow: "hidden",
        display: "flex",
      }}
    >
      {/* Main Dashboard */}
      <div
        style={{
          flex: 1,
          display: "grid",
          gridTemplateColumns: "1fr 1fr",
          gridTemplateRows: "1fr 1fr",
          gap: 16,
          padding: 16,
          boxSizing: "border-box",
          overflow: "hidden",
        }}
      >
        {[0, 1, 2, 3].map((i) => {
          const p = panelSrcs[i % panelSrcs.length];
          return (
            <div
              key={i}
              style={{
                border: "1px solid #222",
                borderRadius: 12,
                overflow: "hidden",
                width: "100%",
                height: "100%",
                position: "relative",
              }}
            >
              <iframe
                src={p.src}
                width="100%"
                height="100%"
                frameBorder="0"
                title={`${p.title} - Panel ${i + 1}`}
                style={{ pointerEvents: "auto" }}
              />
              <button
                onClick={() => triggerAnalysis(p.title)}
                disabled={loading}
                style={{
                  position: "absolute",
                  top: 12,
                  right: 12,
                  width: 40,
                  height: 40,
                  borderRadius: "50%",
                  border: "none",
                  backgroundColor: "#111217",
                  color: "white",
                  fontSize: 18,
                  cursor: loading ? "default" : "pointer",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  boxShadow: "0 2px 8px rgba(0,0,0,0.15)",
                  zIndex: 10,
                  transition: "all 0.2s ease",
                  opacity: loading ? 0.6 : 1,
                }}
                onMouseEnter={(e) => {
                  if (loading) return;
                  e.currentTarget.style.backgroundColor = "#111217";
                  e.currentTarget.style.transform = "scale(1.05)";
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = "#111217";
                  e.currentTarget.style.transform = "scale(1)";
                }}
                title={`Analyze: ${p.title}`}
              >
                🤖
              </button>
            </div>
          );
        })}
      </div>

      {/* Chat Side Panel */}
      <div
        style={{
          position: "fixed",
          left: 0,
          top: 0,
          width: chatOpen ? 350 : 0,
          height: "100vh",
          backgroundColor: "#1a1a1a",
          borderRight: "1px solid #333",
          display: "flex",
          flexDirection: "column",
          transition: "width 0.3s ease",
          zIndex: 1000,
          overflow: "hidden",
        }}
      >
        <div
          style={{
            padding: "16px 12px",
            borderBottom: "1px solid #333",
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            flexShrink: 0,
          }}
        >
          <h3 style={{ margin: 0, color: "#fff", fontSize: 16 }}>
            Metric Analyzer{activePanelTitle ? ` — ${activePanelTitle}` : ""}
          </h3>
          <button
            onClick={() => setChatOpen(false)}
            style={{
              background: "none",
              border: "none",
              color: "#999",
              fontSize: 18,
              cursor: "pointer",
              padding: "4px 8px",
            }}
            title="Close"
          >
            ✕
          </button>
        </div>

        <div
          style={{
            flex: 1,
            overflowY: "auto",
            padding: 12,
            display: "flex",
            flexDirection: "column",
            gap: 12,
          }}
        >
          {messages.map((msg, idx) => (
            <div
              key={idx}
              style={{
                display: "flex",
                justifyContent: msg.role === "user" ? "flex-end" : "flex-start",
              }}
            >
              <div
                style={{
                  maxWidth: "85%",
                  padding: "10px 12px",
                  borderRadius: 8,
                  backgroundColor: msg.role === "user" ? "#4a7c59" : "#333",
                  color: "#fff",
                  fontSize: 14,
                  lineHeight: 1.4,
                  wordWrap: "break-word",
                  whiteSpace: "pre-wrap",
                }}
              >
                {msg.content}
              </div>
            </div>
          ))}

          {loading && (
            <div style={{ display: "flex", justifyContent: "flex-start" }}>
              <div
                style={{
                  padding: "10px 12px",
                  borderRadius: 8,
                  backgroundColor: "#333",
                  color: "#999",
                  fontSize: 14,
                }}
              >
                Thinking...
              </div>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>

        <div
          style={{
            padding: 12,
            borderTop: "1px solid #333",
            display: "flex",
            gap: 8,
            flexShrink: 0,
          }}
        >
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
            placeholder="Ask about the spike..."
            style={{
              flex: 1,
              padding: "8px 12px",
              borderRadius: 6,
              border: "1px solid #444",
              backgroundColor: "#222",
              color: "#fff",
              fontSize: 14,
              outline: "none",
            }}
            disabled={loading}
          />
          <button
            onClick={handleSendMessage}
            disabled={loading}
            style={{
              padding: "8px 12px",
              borderRadius: 6,
              border: "none",
              backgroundColor: loading ? "#555" : "#4a7c59",
              color: "#fff",
              cursor: loading ? "default" : "pointer",
              fontSize: 14,
              fontWeight: 500,
              transition: "background-color 0.2s",
              opacity: loading ? 0.6 : 1,
            }}
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
}