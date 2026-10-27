import React, { useState } from "react";
import { useNavigate } from "react-router";

export default function SellPage() {
  const [title, setTitle] = useState("");
  const [type, setType] = useState("ITEM");
  const [price, setPrice] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const navigate = useNavigate();

  function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");
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
        {error && <p>{error}</p>}
        <button disabled={submitting}>{buttonText}</button>
      </form>
    </div>
  );
}
