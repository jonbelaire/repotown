import "./globals.css";
import Head from "next/head";
import type { Metadata } from "next";

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
      <Head>
        <title>Tramonti Trattoria</title>
        <meta
          name="description"
          content="Authentic Italian dining at Tramonti Trattoria."
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <body>{children}</body>
    </html>
  );
}
