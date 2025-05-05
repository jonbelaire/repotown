"use client";
import Head from "next/head";
import Link from "next/link";
import { useState, useEffect } from "react";
import { FaBars, FaTimes } from "react-icons/fa";

// Menu items data
const menuItemsData = [
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

// Navigation items
const navItems = [
  { href: "#home", label: "Home" },
  { href: "#menu", label: "Menu" },
  { href: "#about", label: "About" },
];

// Modified Header Component
function Header() {
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

  return (
    <header
      style={{
        position: "absolute",
        top: 0,
        left: 0,
        right: 0,
        zIndex: 50,
        color: "#ffffff",
      }}
    >
      <nav style={{ maxWidth: "1280px", margin: "0 auto", padding: "1rem" }}>
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "flex-end",
            height: "5rem",
          }}
        >
          <button
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            style={{
              display: "block",
              color: "#ffffff",
              background: "none",
              border: "none",
              cursor: "pointer",
            }}
          >
            {isMenuOpen ? <FaTimes size={24} /> : <FaBars size={24} />}
          </button>
        </div>
      </nav>
      {isMenuOpen && (
        <div
          style={{
            position: "fixed",
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: "rgba(47, 47, 47, 0.95)",
            zIndex: 50,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
            }}
          >
            <ul
              style={{
                display: "flex",
                flexDirection: "column",
                gap: "2rem",
                listStyle: "none",
                margin: 0,
                padding: 0,
              }}
            >
              {navItems.map((item) => (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    onClick={() => setIsMenuOpen(false)}
                    style={{
                      color:
                        currentHash === item.href.slice(1)
                          ? "#d4a017"
                          : "#ffffff",
                      fontFamily: "'Playfair Display', serif",
                      fontSize: "1.5rem",
                      fontWeight: 600,
                      textTransform: "uppercase",
                      letterSpacing: "0.15em",
                      transition: "all 0.3s ease-in-out",
                    }}
                    onMouseOver={(e) =>
                      (e.currentTarget.style.color = "#d4a017")
                    }
                    onMouseOut={(e) =>
                      (e.currentTarget.style.color =
                        currentHash === item.href.slice(1)
                          ? "#d4a017"
                          : "#ffffff")
                    }
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

// Modified Hero Component
function Hero() {
  return (
    <section
      id="home"
      style={{
        position: "relative",
        background: "url(/hero.png) no-repeat center center/cover",
        height: "100vh",
      }}
    >
      <div
        style={{
          position: "absolute",
          inset: 0,
          background:
            "linear-gradient(to bottom, rgba(0, 0, 0, 0.6), rgba(0, 0, 0, 0.2)",
        }}
      ></div>
      <div
        style={{
          position: "relative",
          zIndex: 10,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          height: "100%",
          textAlign: "center",
          color: "#ffffff",
          padding: "1rem",
        }}
      >
        <img
          src="/logo.svg"
          alt="Tramonti Trattoria"
          style={{
            height: "32rem",
            marginBottom: "2rem",
            filter: "drop-shadow(0 0 10px rgba(0, 0, 0, 0.5))",
          }}
        />
        <p
          style={{
            fontSize: "1.25rem",
            marginBottom: "2rem",
            maxWidth: "32rem",
          }}
        >
          Savor authentic Italian cuisine in a cozy setting.
        </p>
        <Link
          href="#menu"
          style={{
            backgroundColor: "#004d00",
            color: "#ffffff",
            padding: "0.75rem 2rem",
            borderRadius: "0.5rem",
            fontSize: "1.125rem",
            fontWeight: 600,
            textTransform: "uppercase",
            letterSpacing: "0.1em",
            transition: "all 0.3s ease-in-out",
          }}
          onMouseOver={(e) =>
            (e.currentTarget.style.backgroundColor = "#a52a2a")
          }
          onMouseOut={(e) =>
            (e.currentTarget.style.backgroundColor = "#004d00")
          }
        >
          View Menu
        </Link>
      </div>
    </section>
  );
}

// Footer Component - unchanged
function Footer() {
  return (
    <footer
      style={{
        backgroundColor: "#2f2f2f",
        color: "#ffffff",
        padding: "1.5rem 0",
        textAlign: "center",
      }}
    >
      <div>
        <p>&copy; 2025 Tramonti Trattoria. All rights reserved.</p>
        <p>Contact: (123) 456-7890 | info@tramontitrattoria.com</p>
      </div>
    </footer>
  );
}

// Main Page Component
export default function Home() {
  return (
    <div
      style={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}
    >
      <Head>
        <title>Tramonti Trattoria</title>
        <meta
          name="description"
          content="Authentic Italian dining at Tramonti Trattoria."
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main style={{ flexGrow: 1 }}>
        <Hero />

        <section id="welcome" style={{ padding: "3rem 1rem" }}>
          <h2
            style={{
              fontFamily: "'Playfair Display', serif",
              fontSize: "2.25rem",
              fontWeight: 700,
              textAlign: "center",
              marginBottom: "2.5rem",
              color: "#004d00",
              position: "relative",
            }}
          >
            Benvenuti!
            <span
              style={{
                display: "block",
                width: "4rem",
                height: "0.25rem",
                backgroundColor: "#d4a017",
                margin: "0.75rem auto 0",
              }}
            ></span>
          </h2>
          <p
            style={{ textAlign: "center", maxWidth: "32rem", margin: "0 auto" }}
          >
            Experience Italy at Tramonti Trattoria. Browse our menu, read our
            story, or contact us to book a table.
          </p>
        </section>

        <section
          id="menu"
          style={{ padding: "3rem 1rem", backgroundColor: "#f8f1e9" }}
        >
          <h2
            style={{
              fontFamily: "'Playfair Display', serif",
              fontSize: "2.25rem",
              fontWeight: 700,
              textAlign: "center",
              marginBottom: "2.5rem",
              color: "#004d00",
              position: "relative",
            }}
          >
            Our Menu
            <span
              style={{
                display: "block",
                width: "4rem",
                height: "0.25rem",
                backgroundColor: "#d4a017",
                margin: "0.75rem auto 0",
              }}
            ></span>
          </h2>
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(auto-fit, minmax(300px, 1fr))",
              gap: "1.5rem",
              maxWidth: "64rem",
              margin: "0 auto",
            }}
          >
            {menuItemsData.map((item) => (
              <div
                key={item.id}
                style={{
                  padding: "1.5rem",
                  border: "1px solid #e5e7eb",
                  borderRadius: "0.75rem",
                  boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
                  backgroundColor: "#ffffff",
                  transition: "all 0.3s ease-in-out",
                }}
                onMouseOver={(e) =>
                  (e.currentTarget.style.transform = "scale(1.05)")
                }
                onMouseOut={(e) =>
                  (e.currentTarget.style.transform = "scale(1)")
                }
              >
                <img
                  src={item.image}
                  alt={item.name}
                  style={{
                    width: "100%",
                    height: "12rem",
                    objectFit: "cover",
                    borderRadius: "0.5rem 0.5rem 0 0",
                  }}
                />
                <div style={{ padding: "1rem" }}>
                  <h3
                    style={{
                      fontSize: "1.5rem",
                      fontWeight: 600,
                      color: "#004d00",
                    }}
                  >
                    {item.name}
                  </h3>
                  <p style={{ color: "#4b5563", fontStyle: "italic" }}>
                    {item.description}
                  </p>
                  <p
                    style={{
                      fontSize: "1.125rem",
                      fontWeight: 700,
                      marginTop: "0.5rem",
                      color: "#a52a2a",
                    }}
                  >
                    ${item.price.toFixed(2)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </section>

        <section id="about" style={{ padding: "3rem 1rem" }}>
          <h2
            style={{
              fontFamily: "'Playfair Display', serif",
              fontSize: "2.25rem",
              fontWeight: 700,
              textAlign: "center",
              marginBottom: "2.5rem",
              color: "#004d00",
              position: "relative",
            }}
          >
            About Us
            <span
              style={{
                display: "block",
                width: "4rem",
                height: "0.25rem",
                backgroundColor: "#d4a017",
                margin: "0.75rem auto 0",
              }}
            ></span>
          </h2>
          <div style={{ maxWidth: "64rem", margin: "0 auto" }}>
            <div
              style={{
                display: "grid",
                gridTemplateColumns: "repeat(auto-fit, minmax(300px, 1fr))",
                gap: "2rem",
                alignItems: "center",
              }}
            >
              <div>
                <img
                  src="/about-image.jpg"
                  alt="About Tramonti"
                  style={{
                    borderRadius: "0.5rem",
                    boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
                    width: "100%",
                  }}
                />
              </div>
              <div>
                <p style={{ marginBottom: "1rem" }}>
                  Tramonti Trattoria is a family-owned restaurant serving
                  authentic Italian cuisine since 1995. Our recipes have been
                  passed down through generations, bringing the true taste of
                  Italy to your table.
                </p>
                <p>
                  We pride ourselves on using only the freshest ingredients and
                  traditional cooking methods to create an unforgettable dining
                  experience.
                </p>
              </div>
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
