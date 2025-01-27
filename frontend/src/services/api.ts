import axios from "axios";

const apiClient = axios.create({
  baseURL: "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

export const getPantryItems = async () => {
  const response = await apiClient.get("/pantry");
  return response.data;
};

export const createPantryItem = async (item) => {
  const response = await apiClient.post("/pantry", item);
  return response.data;
};

export const updatePantryItem = async (id, item) => {
  const response = await apiClient.put(`/pantry/${id}`, item);
  return response.data;
};

export const deletePantryItem = async (id) => {
  const response = await apiClient.delete(`/pantry/${id}`);
  return response.data;
};
