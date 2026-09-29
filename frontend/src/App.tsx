import { useEffect, useState } from "react";

type Health = {
  status: string;
  database: string;
};

function App() {
  const [health, setHealth] = useState<Health | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/health")
      .then((res) => res.json())
      .then((data: Health) => setHealth(data))
      .catch(() => setError("Could not reach the API"));
  }, []);

  return (
    <main>
      <h1>Seraphim</h1>
      {error && <p>{error}</p>}
      {health && <p>Database: {health.database}</p>}
      {!health && !error && <p>Checking...</p>}
    </main>
  );
}

export default App;
