import { createClient } from "contentful";

export function getContentfulClient() {
  if (!process.env.CONTENTFUL_SPACE_ID) {
    throw new Error("CONTENTFUL_SPACE_ID is required");
  }

  if (!process.env.CONTENTFUL_ACCESS_TOKEN) {
    throw new Error("CONTENTFUL_ACCESS_TOKEN is required");
  }

  return createClient({
    space: process.env.CONTENTFUL_SPACE_ID,
    accessToken: process.env.CONTENTFUL_ACCESS_TOKEN,
  });
}

// Generic fetch function
export async function getEntries<T>(contentType: string): Promise<T[]> {
  const client = getContentfulClient();
  const response = await client.getEntries({ content_type: contentType });
  return response.items.map((item) => item.fields as T);
}
