import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router";
import type { Listing } from "../types";

export default function FeedPage() {
  const [listings, setListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchParams, setSearchParams] = useSearchParams();

  const q = searchParams.get("q") || "";
  const [text, setText] = useState(q);

  const query = searchParams.toString();

  function handleSearch(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    const trimmer = text.trim();
    if (trimmer === "") {
      setSearchParams({});
      return;
    }

    setSearchParams({ q: trimmer });
  }

  useEffect(() => {
    function load(url: string) {
      fetch(url)
        .then((res) => {
          if (!res.ok) {
            return [];
          }
          return res.json();
        })
        .then((data: Listing[]) => {
          setListings(data);
          setLoading(false);
        });
    }

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        load(
          `/api/listings?${query}&lat=${pos.coords.latitude}&lng=${pos.coords.longitude}`,
        );
      },
      () => {
        load(`/api/listings?${query}`);
      },
    );
  }, [query]);

  return (
    <main>
      <div className="feed-header">
        <h1>Listings</h1>
        <form className="search-form" onSubmit={handleSearch}>
          <input
            type="search"
            placeholder="Search listings"
            aria-label="Search listings"
            value={text}
            onChange={(e) => setText(e.target.value)}
          />
        </form>
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

      {!loading && listings.length === 0 && q === "" && (
        <p className="feed-message">
          No listings nearby yet. Be the first to sell something!
        </p>
      )}

      {!loading && listings.length === 0 && q !== "" && (
        <p className="feed-message">No listings match "{q}"</p>
      )}
      <div className="feed-grid">
        {listings.map((listing) => (
          <Link
            key={listing.id}
            to={`/listings/${listing.id}`}
            className="card-link"
          >
            <article className="listing-card">
              {listing.coverUrl && (
                <img
                  className="card-photo"
                  src={listing.coverUrl}
                  alt={listing.title}
                />
              )}
              {!listing.coverUrl && (
                <div className="card-photo card-placeholder">
                  {listing.type === "ITEM" && "📦"}
                  {listing.type === "SERVICE" && "🛠️"}
                </div>
              )}
              <h2>{listing.title}</h2>
              <p className="listing-price">
                {listing.price === 0 && "Free"}
                {listing.price > 0 && `${listing.price} den`}
              </p>
              <p className="listing-type">{listing.type}</p>
            </article>
          </Link>
        ))}
      </div>
    </main>
  );
}
