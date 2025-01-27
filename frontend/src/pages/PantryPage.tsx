import React, { useState, useEffect } from "react";
import { Container, Typography, Button } from "@mui/material"; // Button をインポート
import PantryItemList from "../components/PantryItemList";
import PantryForm from "../components/PantryForm";
import {
  getPantryItems,
  addPantryItem,
  PantryItem,
  getRecipeSuggestion,
} from "../services/api";

const PantryPage: React.FC = () => {
  const [items, setItems] = useState<PantryItem[]>([]);
  const [recipe, setRecipe] = useState<string>(""); // レシピ提案の状態を追加

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      const data = await getPantryItems();
      setItems(data);
    } catch (error) {
      console.error("Error fetching pantry items:", error);
    }
  };

  const handleAddItem = async (item) => {
    try {
      await addPantryItem(item);
      fetchItems();
    } catch (error) {
      console.error("Error adding pantry item:", error);
    }
  };

  const handleSuggestRecipe = async () => {
    // レシピ提案ボタンのクリックハンドラー
    try {
      const data = await getRecipeSuggestion();
      setRecipe(data.recipe); // レシピ提案APIの結果をstateに設定
    } catch (error) {
      console.error("Error suggesting recipe:", error);
    }
  };

  return (
    <Container maxWidth="md" sx={{ mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        Pantry Management
      </Typography>
      <PantryForm onSubmit={handleAddItem} />
      <Button variant="contained" color="primary" onClick={handleSuggestRecipe}>
        {" "}
        // レシピ提案ボタンを追加 レシピ提案
      </Button>
      {recipe && ( // レシピ提案がある場合のみ表示
        <Typography variant="body1" mt={2}>
          レシピ提案: {recipe}
        </Typography>
      )}
      <PantryItemList
        items={items}
        onEdit={(item) => console.log("Edit:", item)}
        onDelete={(id) => console.log("Delete:", id)}
      />
    </Container>
  );
};

export default PantryPage;
