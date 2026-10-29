import React, { useState } from "react";
import { useNavigate } from "react-router";

export default function SellPage() {
  const [title, setTitle] = useState("");
  const [type, setType] = useState("ITEM");
  const [price, setPrice] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [photos, setPhotos] = useState<File[]>([]);

  const navigate = useNavigate();

  function removePhoto(index: number) {
    setPhotos(photos.filter((_, i) => i !== index));
  }

  function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");
    if (photos.length > 5) {
      setError("you can add up to 5 photos");
      return;
    }
    if (photos.length === 0 && type === "ITEM") {
      setError("items need atleast one photo");
      return;
    }
    setSubmitting(true);

    const token = localStorage.getItem("token");

    if (!token) {
      navigate("/login");
      return;
    }

    navigator.geolocation.getCurrentPosition(
      async (position) => {
        const res = await fetch("/api/listings", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            type: type,
            title: title,
            description: description,
            price: Number(price),
            lat: position.coords.latitude,
            lng: position.coords.longitude,
          }),
        });
        const data = await res.json();
        if (res.status === 401) {
          localStorage.removeItem("token");
          navigate("/login");
          return;
        }

        if (!res.ok) {
          setError(data.error);
          setSubmitting(false);
          return;
        }

        for (const photo of photos) {
          const formData = new FormData();
          formData.append("photo", photo);

          const photoRes = await fetch(`/api/listings/${data.id}/photos`, {
            method: "POST",
            headers: {
              Authorization: `Bearer ${token}`,
            },
            body: formData,
          });
          if (!photoRes.ok) {
            const photoData = await photoRes.json();
            setError(photoData.error);
            setSubmitting(false);
            return;
          }
        }

        navigate("/");
      },
      () => {
        setError("We need your location to post a listing");
        setSubmitting(false);
      },
    );
  }

  let itemClass = "type-option";

  if (type === "ITEM") {
    itemClass = "type-option selected";
  }

  let serviceClass = "type-option";

  if (type === "SERVICE") {
    serviceClass = "type-option selected";
  }

  let buttonText = "Submit listing";
  if (submitting) {
    buttonText = "Posting...";
  }

  return (
    <div className="auth-page">
      <form onSubmit={handleSubmit} className="auth-form">
        <h1>Sell something</h1>
        <input
          placeholder="What are you selling?"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
        <div className="type-toggle">
          <button
            type="button"
            className={itemClass}
            onClick={() => setType("ITEM")}
          >
            Item
          </button>
          <button
            type="button"
            className={serviceClass}
            onClick={() => setType("SERVICE")}
          >
            Service
          </button>
        </div>
        <input
          placeholder="Price in MKD"
          type="number"
          value={price}
          onChange={(e) => setPrice(e.target.value)}
        />
        <textarea
          placeholder="Description"
          name="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        ></textarea>
        <input
          type="file"
          accept="image/jpeg,image/png,image/webp"
          multiple
          onChange={(e) => {
            if (e.target.files) {
              const picked = Array.from(e.target.files);

              setPhotos([...photos, ...picked]);
            }
            e.target.value = "";
          }}
        />

        {photos.length > 0 && (
          <ul className="photo-list">
            {photos.map((photo, index) => (
              <li key={index}>
                {photo.name}
                <button type="button" onClick={() => removePhoto(index)}>
                  Remove
                </button>
              </li>
            ))}
          </ul>
        )}
        {error && <p>{error}</p>}
        <button disabled={submitting}>{buttonText}</button>
      </form>
    </div>
  );
}
