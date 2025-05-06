import Header from "./components/header";
import Hero from "./components/hero";
import Welcome from "./components/welcome";
import Menu from "./components/menu";
import About from "./components/about";
import Footer from "./components/footer";

export default function Page() {
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <main className="flex-grow">
        <Hero />
        <Welcome />
        <Menu />
        <About />
      </main>
      <Footer />
    </div>
  );
}
