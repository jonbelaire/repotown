import type { IMenuItem } from "./types";

export const MENU_ITEMS: IMenuItem[] = [
  {
    id: 1,
    name: "Margherita Pizza",
    description: "Classic pizza with tomato, mozzarella, and basil.",
    price: 12.99,
    image: "/pizza.jpg",
  },
  {
    id: 2,
    name: "Spaghetti Carbonara",
    description: "Pasta with egg, pecorino, guanciale, and black pepper.",
    price: 14.99,
    image: "/spaghetti.jpg",
  },
  {
    id: 3,
    name: "Tiramisu",
    description: "Coffee-flavored Italian dessert.",
    price: 6.99,
    image: "/tiramisu.jpg",
  },
];
