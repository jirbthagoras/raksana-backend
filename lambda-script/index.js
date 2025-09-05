import { GoogleGenAI } from "@google/genai";
import { Pool } from "pg";

const pool = new Pool({
  host: process.env.DB_HOST,
  port: Number(process.env.DB_PORT),
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
});

export const handler = async () => {
  try {
    const ai = new GoogleGenAI({ apiKey: process.env.GEMINI_API_KEY });

    const prompt = process.env.SYSTEM_INSTRUCTION;

    const response = await ai.models.generateContent({
      model: "gemini-2.0-flash-001",
      contents: [prompt],
    });

    let challengeText = response.candidates[0].content.parts[0].text;

    challengeText = challengeText.replace(/```json\s*|```/g, "").trim();

    const challenge = JSON.parse(challengeText);

    const detailQuery = `
      INSERT INTO details (name, description, point_gain, created_at, updated_at)
      VALUES ($1, $2, $3, NOW(), NOW())
      RETURNING id
    `;
    const detailValues = [
      challenge.title,
      challenge.description,
      challenge.points,
    ];
    const detailRes = await pool.query(detailQuery, detailValues);
    const detailId = detailRes.rows[0].id;

    const getDayQuery = `
      SELECT COALESCE(MAX(day), 0) + 1 AS next_day
      FROM challenges
    `;

    const result = await pool.query(getDayQuery);
    const day = result.rows[0].next_day;

    console.log(challenge);

    const difficulty = () => {
      if (challenge.points >= 50 && challenge.points <= 100) {
        return "easy";
      } else if (challenge.points > 250) {
        return "hard";
      } else {
        return "normal";
      }
    };

    const challengeQuery = `
      INSERT INTO challenges (detail_id, day, difficulty)
      VALUES ($1, $2, $3)
      RETURNING *
    `;
    const challengeValues = [detailId, day, difficulty()];
    const challengeRes = await pool.query(challengeQuery, challengeValues);

    return {
      statusCode: 200,
      body: JSON.stringify({
        detail: challenge,
        challengeRow: challengeRes.rows[0],
      }),
    };
  } catch (error) {
    console.error(error);
    return {
      statusCode: 500,
      body: JSON.stringify({ error: error.message }),
    };
  }
};
