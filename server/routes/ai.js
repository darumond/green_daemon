import express from "express";
import { GoogleGenerativeAI } from "@google/generative-ai";

const router = express.Router();

router.post("/query", async (req, res) => {
  try {
    const { query, context } = req.body ?? {};

    if (!query || typeof query !== "string") {
      return res.status(400).json({ message: "Missing 'query' (string)." });
    }

    const apiKey = process.env.GEMINI_API_KEY;
    if (!apiKey) {
      return res.status(500).json({ message: "GEMINI_API_KEY not set on server." });
    }

    const genAI = new GoogleGenerativeAI(apiKey);

    // If this model name errors, swap to another Gemini model name you have enabled.
    const model = genAI.getGenerativeModel({ model: "gemini-2.5-flash" });

    const prompt = `
You are an SRE assistant analyzing Grafana dashboards and performance spikes.
Context: ${context ?? "N/A"}

User question: ${query}

Return:
- short diagnosis (1-2 sentences)
- 3 likely causes (bullets)
- 3 checks to run (bullets)
- 1 mitigation (1 sentence)
`.trim();

    const result = await model.generateContent(prompt);
    const text = result.response.text();

    return res.json({ message: text });
  } catch (err) {
    console.error(err);
    return res.status(500).json({ message: "Gemini request failed." });
  }
});

export default router;