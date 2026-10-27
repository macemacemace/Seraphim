import { useEffect, useState } from "react";
import { Link } from "react-router";

type Listing = {
  id: string;
  title: string;
  price: number;
  sellerId: string;
  type: string;
  description: string;
  lat: number;
  lng: number;
};

export default function FeedPage() {
  const [listings, setListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    function load(url: string) {
      fetch(url)
        .then((res) => res.json())
        .then((data: Listing[]) => {
          setListings(data);
          setLoading(false);
        });
    }

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        load(
          `/api/listings?lat=${pos.coords.latitude}&lng=${pos.coords.longitude}`,
        );
      },
      () => {
        load(`/api/listings`);
      },
    );
  }, []);

  return (
    <main>
      <div className="feed-header">
        <h1>Listings</h1>
        <Link to="/sell" className="sell-button">
          Sell something
        </Link>
      </div>
      {loading && (
        <div className="feed-loading">
          <div className="dots">
            <span></span>
            <span></span>
            <span></span>
          </div>
          <p className="feed-message">Finding listings near you...</p>
        </div>
      )}

      {!loading && listings.length === 0 && (
        <p className="feed-message">
          No listings nearby yet. Be the first to sell something!
        </p>
      )}
      <div className="feed-grid">
        {listings.map((listing) => (
          <article key={listing.id} className="listing-card">
            <h2>{listing.title}</h2>
            <p className="listing-price">
              {listing.price == 0 && "Free"}
              {listing.price > 0 && `${listing.price} den`}
            </p>
            <p className="listing-type">{listing.type}</p>
          </article>
        ))}
      </div>
    </main>
  );
}
