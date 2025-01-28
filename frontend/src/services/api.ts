import axios from "axios";

interface Food {
  id: number;
  name: string;
  description: string;
  quantity: number;
  unit: string;
  expiry_date: string;
}

const apiClient = axios.create({
  baseURL: "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

export const getPantryItems = async (): Promise<Food[]> => {
  const response = await apiClient.get("/pantry");
  return response.data;
};

export const createPantryItem = async (item: Food): Promise<Food> => {
  const response = await apiClient.post("/pantry", item);
  return response.data;
};

export const updatePantryItem = async (
  id: number,
  item: Food
): Promise<Food> => {
  const response = await apiClient.put(`/pantry/${id}`, item);
  return response.data;
};

export const deletePantryItem = async (id: number): Promise<void> => {
  await apiClient.delete(`/pantry/${id}`);
};
