import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router";
import type { Listing, Photo } from "../types";
import { Link } from "react-router";
import { useCurrentUser } from "../useCurrentUser";
import LoadingDots from "../components/LoadingDots";

export default function ListingPage() {
  const { id } = useParams();
  const [listing, setListing] = useState<Listing | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [selected, setSelected] = useState(0);
  const user = useCurrentUser();
  const navigate = useNavigate();

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
        <LoadingDots />
      </main>
    );
  }
  let icon = "📦";
  if (listing.type === "SERVICE") {
    icon = "🛠️";
  }

  let photos: Photo[] = [];
  if (listing.photos) {
    photos = listing.photos;
  }

  function showNext() {
    if (selected === photos.length - 1) {
      setSelected(0);
    } else {
      setSelected(selected + 1);
    }
  }

  function showPrevious() {
    if (selected === 0) {
      setSelected(photos.length - 1);
    } else {
      setSelected(selected - 1);
    }
  }

  const listedDate = new Date(listing.createdAt).toLocaleDateString("en-GB", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });

  function toggleStatus() {
    if (listing === null) {
      return;
    }

    let newStatus = "SOLD";
    if (listing.status === "SOLD") {
      newStatus = "ACTIVE";
    }

    fetch(`/api/listings/${listing.id}/status`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify({ status: newStatus }),
    }).then((res) => {
      if (res.ok) {
        setListing({ ...listing, status: newStatus });
      }
    });
  }

  function handleDelete() {
    if (listing === null) {
      return;
    }

    const sure = window.confirm("Delete this listing? This can't be undone.");
    if (!sure) {
      return;
    }
    fetch(`/api/listings/${listing.id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    }).then((res) => {
      if (res.ok) {
        navigate("/");
      }
    });
  }

  let statusButtonText = "Mark as sold";
  if (listing.status === "SOLD") {
    statusButtonText = "Mark as active";
  }

  return (
    <main className="listing-page">
      <Link to="/" className="back-link">
        ← Back to listings
      </Link>
      {photos.length === 0 && <div className="listing-photo">{icon}</div>}
      {photos.length > 0 && (
        <>
          <div className="photo-viewer">
            <img
              className="listing-photo"
              src={photos[selected].url}
              alt={listing.title}
            />

            {photos.length > 1 && (
              <>
                <button
                  className="photo-arrow prev"
                  type="button"
                  onClick={showPrevious}
                >
                  ‹
                </button>
                <button
                  className="photo-arrow next"
                  type="button"
                  onClick={showNext}
                >
                  ›
                </button>
              </>
            )}
          </div>

          <div className="thumbnails">
            {photos.map((photo, index) => {
              let thumbClass = "thumbnail";
              if (index === selected) {
                thumbClass = "thumbnail selected";
              }

              return (
                <button
                  key={photo.id}
                  type="button"
                  className={thumbClass}
                  onClick={() => setSelected(index)}
                >
                  <img
                    src={photo.url}
                    alt={`${listing.title} photo ${index + 1}`}
                  />
                </button>
              );
            })}
          </div>
        </>
      )}

      {listing.status === "SOLD" && <span className="sold-badge">Sold</span>}

      <h1>{listing.title}</h1>
      <p className="listing-price">
        {listing.price === 0 && "Free"}
        {listing.price > 0 && `${listing.price} den`}
      </p>
      <p className="listing-type">{listing.type}</p>
      <p className="listing-description">{listing.description}</p>
      <p className="listing-date">Listed on {listedDate}</p>

      {user && user.id === listing.sellerId && (
        <div className="seller-actions">
          <>
            <button
              type="button"
              className="status-button"
              onClick={toggleStatus}
            >
              {statusButtonText}
            </button>
            <button
              type="button"
              className="delete-button"
              onClick={handleDelete}
            >
              Delete listing
            </button>
          </>
        </div>
      )}
    </main>
  );
}
