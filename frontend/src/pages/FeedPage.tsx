import { useEffect, useState } from "react";

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

  useEffect(() => {
    function load(url: string) {
      fetch(url)
        .then((res) => res.json())
        .then((data: Listing[]) => setListings(data));
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
      <h1>Listings</h1>
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
