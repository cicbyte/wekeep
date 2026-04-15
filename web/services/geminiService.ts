
import { GoogleGenAI, Type } from "@google/genai";
import { ParseResult } from "../types";

const parseArticleContent = async (text: string): Promise<ParseResult> => {
  if (!process.env.API_KEY) {
    throw new Error("API Key is missing");
  }

  const ai = new GoogleGenAI({ apiKey: process.env.API_KEY });

  const prompt = `
    Analyze the following text, which is likely copied from a WeChat Official Account article or a similar blog post.
    Extract the following information:
    1. Title of the article.
    2. Author Name (Official Account Name).
    3. A brief summary (1-2 sentences).
    4. Relevant tags (max 3).
    5. The URL if present in the text.
    6. The FULL article content, converted into clean, well-formatted Markdown. Remove advertisements, "read more" links, and header/footer fluff. Preserve headers, lists, and main content.

    If the author is missing, guess it based on context or use "Unknown Author".
    If the title is missing, generate a descriptive title.

    Input Text:
    """
    ${text}
    """
  `;

  try {
    const response = await ai.models.generateContent({
      model: "gemini-3-flash-preview",
      contents: prompt,
      config: {
        responseMimeType: "application/json",
        responseSchema: {
          type: Type.OBJECT,
          properties: {
            title: { type: Type.STRING, description: "The article title" },
            author: { type: Type.STRING, description: "The official account or author name" },
            summary: { type: Type.STRING, description: "A concise summary" },
            content: { type: Type.STRING, description: "The full article content in Markdown format" },
            tags: { 
              type: Type.ARRAY, 
              items: { type: Type.STRING },
              description: "Up to 3 relevant tags" 
            },
            url: { type: Type.STRING, description: "The article URL if found, else empty string" }
          },
          required: ["title", "author", "summary", "content", "tags"],
        },
      },
    });

    const result = response.text;
    if (!result) throw new Error("No response from AI");
    
    return JSON.parse(result) as ParseResult;

  } catch (error) {
    console.error("Gemini Parse Error:", error);
    throw new Error("Failed to parse content using AI.");
  }
};

export { parseArticleContent };
