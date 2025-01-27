import axios from "axios";

export interface PantryItem {
  ID: number;
  Name: string;
  Quantity: number;
}

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || "http://localhost:8080/api",
});

export const getPantryItems = async (): Promise<PantryItem[]> => {
  const response = await api.get("/pantry");
  return response.data;
};

export const addPantryItem = async (
  item: Omit<PantryItem, "ID">
): Promise<PantryItem> => {
  const response = await api.post("/pantry", item);
  return response.data;
};

export const updatePantryItem = async (
  id: number,
  item: Partial<PantryItem>
): Promise<PantryItem> => {
  const response = await api.put(`/pantry/${id}`, item);
  return response.data;
};

export const deletePantryItem = async (id: number): Promise<void> => {
  await api.delete(`/pantry/${id}`);
};

export const getRecipeSuggestion = async () => {
  // レシピ提案APIを呼び出す関数を追加
  const response = await api.get("/pantry/recipe");
  return response.data;
};
