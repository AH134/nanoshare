export function Loading() {
  return (
    <div className="hero bg-base-200 min-h-screen">
      <div className="hero-content text-center">
        <div className="max-w-md">
          <h1 className="text-5xl font-extrabold">NANOSHARE</h1>
          <p className="py-6 flex items-center justify-center gap-2">
            Loading Nanoshare...
            <span className="loading loading-spinner loading-xs"></span>
          </p>
        </div>
      </div>
    </div>
  );
}
