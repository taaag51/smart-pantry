import React, { useState } from "react";
import FoodForm from "../components/FoodForm";
import FoodList from "../components/FoodList";

const FoodPage: React.FC = () => {
  const [refreshKey, setRefreshKey] = useState(0);

  const handleFoodCreated = () => {
    setRefreshKey((prev) => prev + 1);
  };

  return (
    <div>
      <h1>Food Management</h1>
      <div style={{ display: "flex", gap: "2rem" }}>
        <div style={{ flex: 1 }}>
          <h2>Add New Food</h2>
          <FoodForm onSubmit={handleFoodCreated} />
        </div>
        <div style={{ flex: 2 }}>
          <h2>Food Items</h2>
          <FoodList key={refreshKey} onFoodDeleted={handleFoodCreated} />
        </div>
      </div>
    </div>
  );
};

export default FoodPage;
