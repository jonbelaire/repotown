import Link from "next/link";

export default function Hero() {
  return (
    <section
      id="home"
      className="relative bg-[url(/hero.png)] bg-no-repeat bg-center bg-cover h-screen"
    >
      <div className="absolute inset-0 bg-gradient-to-b from-[rgba(0,0,0,0.6)] to-[rgba(0,0,0,0.2)]"></div>
      <div className="relative z-10 flex flex-col items-center justify-center h-full text-center text-white px-4">
        <img
          src="/logo.svg"
          alt="Tramonti Trattoria"
          className="max-w-full h-auto xl:w-[32rem] drop-shadow-lg"
        />
        <p className="text-xl mb-8 max-w-2xl">
          Savor authentic Italian cuisine in a cozy setting.
        </p>
        <Link
          href="#menu"
          className="bg-green-900 text-white py-3 px-8 rounded-lg text-lg font-semibold uppercase tracking-wider transition-colors duration-300 hover:bg-red-800"
        >
          View Menu
        </Link>
      </div>
    </section>
  );
}
