"use client";
import { useEffect, useState } from "react";
import { FaBars, FaTimes } from "react-icons/fa";
import Link from "next/link";

export default function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [currentHash, setCurrentHash] = useState("home");

  useEffect(() => {
    const updateHash = () => {
      setCurrentHash(window.location.hash.slice(1) || "home");
    };
    window.addEventListener("hashchange", updateHash);
    updateHash();
    return () => window.removeEventListener("hashchange", updateHash);
  }, []);

  const navItems = [
    { href: "#home", label: "Home" },
    { href: "#menu", label: "Menu" },
    { href: "#about", label: "About" },
  ];

  return (
    <header className="absolute top-0 left-0 right-0 z-50 text-white">
      <nav className="max-w-7xl mx-auto px-4">
        <div className="flex items-center justify-end h-20">
          <button
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            className="text-white bg-transparent border-none cursor-pointer z-[60]"
          >
            {isMenuOpen ? <FaTimes size={24} /> : <FaBars size={24} />}
          </button>
        </div>
      </nav>
      {isMenuOpen && (
        <div className="fixed inset-0 bg-[rgba(47,47,47,0.95)] z-50 flex items-center justify-center">
          <div className="flex flex-col items-center justify-center h-full">
            <ul className="flex flex-col gap-8 list-none m-0 p-0">
              {navItems.map((item) => (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    onClick={() => setIsMenuOpen(false)}
                    className={`font-['Playfair_Display'] text-2xl font-semibold uppercase tracking-widest transition-colors duration-300 ${
                      currentHash === item.href.slice(1)
                        ? "text-[#d4a017]"
                        : "text-white hover:text-[#d4a017]"
                    }`}
                  >
                    {item.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>
      )}
    </header>
  );
}
