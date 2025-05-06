import type { IMenuItem } from "./types";

export default function MenuItem({ item }: { item: IMenuItem }) {
  return (
    <div className="p-6 border border-gray-200 rounded-xl shadow-lg bg-white transition-transform duration-300 ease-in-out hover:scale-105">
      <img
        src={item.image}
        alt={item.name}
        className="w-full h-48 object-cover rounded-t-lg"
      />
      <div className="p-4">
        <h3 className="text-2xl font-semibold text-green-800">{item.name}</h3>
        <p className="text-gray-600 italic">{item.description}</p>
        <p className="text-lg font-bold mt-2 text-red-800">
          ${item.price.toFixed(2)}
        </p>
      </div>
    </div>
  );
}
