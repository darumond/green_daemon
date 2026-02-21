import express from "express";
import cors from "cors";
import dotenv from "dotenv";
import aiRoutes from "./routes/ai.js";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

// If frontend and backend are on same origin, you can remove cors() entirely.
// Keeping it here is fine for local dev.
app.use(cors());
app.use(express.json());

app.use("/api/ai", aiRoutes);

app.get("/health", (req, res) => res.json({ ok: true }));

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});