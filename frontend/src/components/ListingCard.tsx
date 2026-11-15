import { Link } from "react-router";
import type { Listing } from "../types";
import { timeAgo } from "../timeAgo";

type ListingCardProps = {
  listing: Listing;
};

export default function ListingCard({ listing }: ListingCardProps) {
  return (
    <Link to={`/listings/${listing.id}`} className="card-link">
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
        {listing.status === "SOLD" && <span className="sold-badge">Sold</span>}
        <p className="listing-price">
          {listing.price === 0 && "Free"}
          {listing.price > 0 && `${listing.price} den`}
        </p>
        <p className="listing-type">{listing.type}</p>
        <p className="listing-date">{timeAgo(listing.createdAt)}</p>
      </article>
    </Link>
  );
}
