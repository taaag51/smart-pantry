import React, { useState } from "react";
import { TextField, Button, Box, Typography, Stack } from "@mui/material";

interface PantryFormProps {
  onSubmit: (item: { Name: string; Quantity: number }) => void;
}

const PantryForm: React.FC<PantryFormProps> = ({ onSubmit }) => {
  const [name, setName] = useState("");
  const [quantity, setQuantity] = useState(1);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError("");

    try {
      await onSubmit({
        Name: name,
        Quantity: quantity,
      });
      setName("");
      setQuantity(1);
    } catch (err) {
      setError("Failed to save item. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4, width: "100%" }}>
      <Typography variant="h6" gutterBottom>
        Add New Item
      </Typography>
      <Stack
        spacing={2}
        direction={{ xs: "column", sm: "row" }}
        sx={{
          mb: 2,
          width: "100%",
          "& .MuiTextField-root": {
            width: { xs: "100%", sm: "auto" },
          },
        }}
      >
        <TextField
          label="Item Name"
          variant="outlined"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          fullWidth
        />
        <TextField
          label="Quantity"
          type="number"
          variant="outlined"
          value={quantity}
          onChange={(e) => setQuantity(Number(e.target.value))}
          required
          inputProps={{ min: 1 }}
          sx={{ maxWidth: 120 }}
        />
        <Button
          type="submit"
          variant="contained"
          size="large"
          disabled={isSubmitting}
          sx={{
            width: { xs: "100%", sm: "auto" },
            mt: { xs: 2, sm: 0 },
          }}
        >
          {isSubmitting ? "Adding..." : "Add Item"}
        </Button>
      </Stack>
    </Box>
  );
};

export default PantryForm;
