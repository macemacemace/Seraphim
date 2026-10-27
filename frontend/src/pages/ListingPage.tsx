import { useEffect, useState } from "react";
import { useParams } from "react-router";
import type { Listing } from "../types";
import { Link } from "react-router";

export default function ListingPage() {
  const { id } = useParams();
  const [listing, setListing] = useState<Listing | null>(null);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    fetch(`/api/listings/${id}`)
      .then((res) => {
        if (res.status === 404) {
          setNotFound(true);
          return null;
        }
        return res.json();
      })
      .then((data) => {
        if (data) {
          setListing(data);
        }
      });
  }, [id]);

  if (notFound) {
    return (
      <main>
        <p>Listing not found</p>
      </main>
    );
  }
  if (listing === null) {
    return (
      <main>
        <p>Loading...</p>
      </main>
    );
  }
  let icon = "📦";
  if (listing.type === "SERVICE") {
    icon = "🛠️";
  }
  return (
    <main className="listing-page">
      <Link to="/" className="back-link">
        ← Back to listings
      </Link>
      <div className="listing-photo">{icon}</div>

      <h1>{listing.title}</h1>
      <p className="listing-price">
        {listing.price === 0 && "Free"}
        {listing.price > 0 && `${listing.price} den`}
      </p>
      <p className="listing-type">{listing.type}</p>
      <p className="listing-description">{listing.description}</p>
    </main>
  );
}
