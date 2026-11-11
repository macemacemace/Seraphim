import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router";
import type { Listing } from "../types";
import { timeAgo } from "../timeAgo";

export default function FeedPage() {
  const [listings, setListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchParams, setSearchParams] = useSearchParams();

  const q = searchParams.get("q") || "";
  const type = searchParams.get("type") || "";
  const [text, setText] = useState(q);

  const [minPrice, setMinPrice] = useState(searchParams.get("minPrice") || "");
  const [maxPrice, setMaxPrice] = useState(searchParams.get("maxPrice") || "");

  const query = searchParams.toString();

  const sort = searchParams.get("sort") || "";

  function handleSearch(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    const next = new URLSearchParams(searchParams);
    const trimmer = text.trim();

    if (trimmer === "") {
      next.delete("q");
    } else {
      next.set("q", trimmer);
    }

    setSearchParams(next);
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

  function changeType(value: string) {
    const next = new URLSearchParams(searchParams);

    if (value === "") {
      next.delete("type");
    } else {
      next.set("type", value);
    }
    setSearchParams(next);
  }

  function typeClass(value: string) {
    if (value === type) {
      return "filter-option selected";
    }
    return "filter-option";
  }

  function changeSort(value: string) {
    const next = new URLSearchParams(searchParams);

    if (value === "") {
      next.delete("sort");
    } else {
      next.set("sort", value);
    }

    setSearchParams(next);
  }

  function handlePrice(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    const next = new URLSearchParams(searchParams);

    if (minPrice === "") {
      next.delete("minPrice");
    } else {
      next.set("minPrice", minPrice);
    }

    if (maxPrice === "") {
      next.delete("maxPrice");
    } else {
      next.set("maxPrice", maxPrice);
    }

    setSearchParams(next);
  }

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
      <div className="filter-row">
        <button
          type="button"
          className={typeClass("")}
          onClick={() => changeType("")}
        >
          All
        </button>
        <button
          type="button"
          className={typeClass("ITEM")}
          onClick={() => changeType("")}
        >
          Item
        </button>
        <button
          type="button"
          className={typeClass("SERVICE")}
          onClick={() => changeType("")}
        >
          Service
        </button>
      </div>
      <form className="price-filter" onSubmit={handlePrice}>
        <input
          type="number"
          min="0"
          placeholder="Min price"
          aria-label="Minimum price"
          value={minPrice}
          onChange={(e) => setMinPrice(e.target.value)}
        />

        <input
          type="number"
          min="0"
          placeholder="Max price"
          aria-label="Maximum price"
          value={maxPrice}
          onChange={(e) => setMaxPrice(e.target.value)}
        />

        <button type="submit">Apply</button>
      </form>

      <select
        aria-label="Sort listings"
        value={sort}
        onChange={(e) => changeSort(e.target.value)}
      >
        <option value="">Newest first</option>
        <option value="oldest">Oldest first</option>
        <option value="price_desc">Price: high to low</option>
        <option value="price_asc">Price: low to high</option>
      </select>
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
              <p className="listing-date">{timeAgo(listing.createdAt)}</p>
            </article>
          </Link>
        ))}
      </div>
    </main>
  );
}
