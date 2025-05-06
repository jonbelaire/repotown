export default function About() {
  return (
    <section id="about" className="py-12 px-4">
      <h2 className="font-['Playfair_Display'] text-4xl font-bold text-center mb-10 text-green-800 relative">
        About Us
        <span className="block w-16 h-1 bg-[#d4a017] mt-3 mx-auto"></span>
      </h2>
      <div className="max-w-5xl mx-auto">
        <div className="grid grid-cols-[repeat(auto-fit,minmax(300px,1fr))] gap-8 items-center">
          <div>
            <img
              src="/about.jpg"
              alt="About Tramonti"
              className="rounded-lg shadow-lg w-full"
            />
          </div>
          <div>
            <p className="mb-4">
              Tramonti Trattoria is a family-owned restaurant serving authentic
              Italian cuisine since 1995. Our recipes have been passed down
              through generations, bringing the true taste of Italy to your
              table.
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
  );
}
