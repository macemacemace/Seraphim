import React, { useState } from "react";
import { useNavigate } from "react-router";

export default function SellPage() {
  const [title, setTitle] = useState("");
  const [type, setType] = useState("ITEM");
  const [price, setPrice] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");

  const navigate = useNavigate();

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

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
        if (!res.ok) {
          setError(data.error);
          return;
        }

        navigate("/");
      },
      () => setError("We need your location to post a listing"),
    );
  }

  return (
    <form onSubmit={handleSubmit}>
      <h1>Sell something</h1>
      <input
        placeholder="What are you selling?"
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <select value={type} onChange={(e) => setType(e.target.value)}>
        <option value="ITEM">Item</option>
        <option value="SERVICE">Service</option>
      </select>
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
      <button>Submit listing</button>
    </form>
  );
}
