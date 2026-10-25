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
      <div>
        {listings.map((listing) => (
          <p key={listing.id}>
            {listing.title}: {listing.price} den
          </p>
        ))}
      </div>
    </main>
  );
}
