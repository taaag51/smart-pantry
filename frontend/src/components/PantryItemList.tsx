import React from "react";
import {
  List,
  ListItem,
  ListItemText,
  IconButton,
  ListItemSecondaryAction,
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";

interface PantryItem {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  Name: string;
  Quantity: number;
  ExpiryDate: string | null;
}

interface PantryItemListProps {
  items: PantryItem[];
  onDeleteItem: (id: number) => void;
  onUpdateItem: (id: number, updatedItem: PantryItem) => void;
}

const PantryItemList: React.FC<PantryItemListProps> = ({
  items,
  onDeleteItem,
  onUpdateItem,
}) => {
  return (
    <List>
      {items.map((item) => (
        <ListItem key={item.ID}>
          <ListItemText
            primary={item.Name}
            secondary={`Quantity: ${item.Quantity}, Expires: ${
              item.ExpiryDate || "N/A"
            }`}
          />
          <ListItemSecondaryAction>
            <IconButton
              edge="end"
              aria-label="edit"
              onClick={() => onUpdateItem(item.ID, item)}
            >
              <EditIcon />
            </IconButton>
            <IconButton
              edge="end"
              aria-label="delete"
              onClick={() => onDeleteItem(item.ID)}
            >
              <DeleteIcon />
            </IconButton>
          </ListItemSecondaryAction>
        </ListItem>
      ))}
    </List>
  );
};

export default PantryItemList;
