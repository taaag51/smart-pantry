import React, { useState, useEffect } from "react";
import axios from "axios";
import PantryForm from "../components/PantryForm";
import PantryItemList from "../components/PantryItemList";
import { Container, Typography, Snackbar, Alert } from "@mui/material";

interface PantryItem {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  Name: string;
  Quantity: number;
  ExpiryDate: string | null;
}

const PantryPage: React.FC = () => {
  const [items, setItems] = useState<PantryItem[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      const response = await axios.get("/api/pantry");
      setItems(response.data);
    } catch (err: any) {
      setError("Failed to fetch pantry items.");
    }
  };

  const handleCreateItem = async (newItem: {
    Name: string;
    Quantity: number;
    ExpiryDate: string | null;
  }) => {
    try {
      const response = await axios.post("/api/pantry", newItem);
      setItems([...items, response.data]);
      setSuccess("Item created successfully.");
      setError(null);
      fetchItems();
    } catch (err: any) {
      setError("Failed to create pantry item.");
      setSuccess(null);
    }
  };

  const handleUpdateItem = async (id: number, updatedItem: PantryItem) => {
    try {
      await axios.put(`/api/pantry/${id}`, updatedItem);
      const updatedItems = items.map((item) =>
        item.ID === id ? updatedItem : item
      );
      setItems(updatedItems);
      setSuccess("Item updated successfully.");
      setError(null);
      fetchItems(); // Refresh item list
    } catch (err: any) {
      setError("Failed to update pantry item.");
      setSuccess(null);
    }
  };

  const handleDeleteItem = async (id: number) => {
    try {
      await axios.delete(`/api/pantry/${id}`);
      const filteredItems = items.filter((item) => item.ID !== id);
      setItems(filteredItems);
      setSuccess("Item deleted successfully.");
      setError(null);
      fetchItems(); // Refresh item list
    } catch (err: any) {
      setError("Failed to delete pantry item.");
      setSuccess(null);
    }
  };

  const handleSnackbarClose = () => {
    setError(null);
    setSuccess(null);
  };

  return (
    <Container maxWidth="md">
      <Typography variant="h4" component="h1" gutterBottom>
        Smart Pantry
      </Typography>

      <PantryForm onCreateItem={handleCreateItem} />

      <PantryItemList
        items={items}
        onUpdateItem={handleUpdateItem}
        onDeleteItem={handleDeleteItem}
      />

      <Snackbar
        open={!!error}
        autoHideDuration={6000}
        onClose={handleSnackbarClose}
      >
        <Alert
          onClose={handleSnackbarClose}
          severity="error"
          sx={{ width: "100%" }}
        >
          {error}
        </Alert>
      </Snackbar>

      <Snackbar
        open={!!success}
        autoHideDuration={6000}
        onClose={handleSnackbarClose}
      >
        <Alert
          onClose={handleSnackbarClose}
          severity="success"
          sx={{ width: "100%" }}
        >
          {success}
        </Alert>
      </Snackbar>
    </Container>
  );
};

export default PantryPage;
