import React, { useState, useEffect } from "react";
import { getPantryItems, deletePantryItem } from "../services/api";

interface Food {
  id: number;
  name: string;
  description: string;
  quantity: number;
  unit: string;
  expiry_date: string;
}

interface FoodListProps {
  onFoodDeleted: () => void;
}

const FoodList: React.FC<FoodListProps> = ({ onFoodDeleted }) => {
  const [foods, setFoods] = useState<Food[]>([]);

  useEffect(() => {
    loadFoods();
  }, [onFoodDeleted]);

  const loadFoods = async () => {
    const data = await getPantryItems();
    setFoods(data);
  };

  const handleDelete = async (id: number) => {
    await deletePantryItem(id);
    onFoodDeleted();
  };

  return (
    <div>
      <h2>Food List</h2>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Description</th>
            <th>Quantity</th>
            <th>Unit</th>
            <th>Expiry Date</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {foods.map((food) => (
            <tr key={food.id}>
              <td>{food.name}</td>
              <td>{food.description}</td>
              <td>{food.quantity}</td>
              <td>{food.unit}</td>
              <td>{food.expiry_date}</td>
              <td>
                <button onClick={() => handleDelete(food.id)}>Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default FoodList;
