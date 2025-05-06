import { MENU_ITEMS } from "./data";
import MenuItem from "./menu-item";

export default function Menu() {
  return (
    <section id="menu" className="py-12 px-4">
      <h2 className="font-['Playfair_Display'] text-4xl font-bold text-center mb-10 text-green-800 relative">
        Our Menu
        <span className="block w-16 h-1 bg-[#d4a017] mt-3 mx-auto"></span>
      </h2>
      <div className="grid grid-cols-[repeat(auto-fit,minmax(300px,1fr))] gap-6 max-w-5xl mx-auto">
        {MENU_ITEMS.map((item) => (
          <MenuItem key={item.id} item={item} />
        ))}
      </div>
    </section>
  );
}
