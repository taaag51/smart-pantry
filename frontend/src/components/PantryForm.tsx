import React, { useState } from "react";
import { TextField, Button, Paper, Grid, Typography } from "@mui/material";

interface PantryFormProps {
  onCreateItem: (newItem: {
    Name: string;
    Quantity: number;
    ExpiryDate: string | null;
  }) => void;
}

const PantryForm: React.FC<PantryFormProps> = ({ onCreateItem }) => {
  const [name, setName] = useState("");
  const [quantity, setQuantity] = useState<number>(1);
  const [expiryDate, setExpiryDate] = useState<string | null>(null);

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    onCreateItem({ Name: name, Quantity: quantity, ExpiryDate: expiryDate });
    setName("");
    setQuantity(1);
    setExpiryDate(null);
  };

  return (
    <Paper elevation={3} style={{ padding: "20px", marginBottom: "20px" }}>
      <Typography variant="h6" component="h2" gutterBottom>
        Add New Item
      </Typography>
      <form onSubmit={handleSubmit}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Quantity"
              type="number"
              value={quantity}
              onChange={(e) => setQuantity(Number(e.target.value))}
              inputProps={{ min: "1" }}
              required
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Expiry Date"
              type="date"
              value={expiryDate || ""}
              onChange={(e) => setExpiryDate(e.target.value)}
              InputLabelProps={{ shrink: true }}
            />
          </Grid>
          <Grid item xs={12}>
            <Button type="submit" variant="contained" color="primary">
              Add Item
            </Button>
          </Grid>
        </Grid>
      </form>
    </Paper>
  );
};

export default PantryForm;
