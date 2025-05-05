import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Tramonti Trattoria – Authentic Italian Cuisine",
  description:
    "Indulge in handcrafted pasta, wood-fired pizza, and family recipes at Tramonti Trattoria since 1995.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="scroll-smooth antialiased bg-[#FAF5EE]">
      <body>{children}</body>
    </html>
  );
}
