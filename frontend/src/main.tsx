import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { BrowserRouter, Routes, Route } from "react-router";
import LoginPage from "./pages/LoginPage.tsx";
import RegisterPage from "./pages/RegisterPage.tsx";
import FeedPage from "./pages/FeedPage.tsx";
import SellPage from "./pages/SellPage.tsx";
import ListingPage from "./pages/ListingPage.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/" element={<FeedPage />} />
        <Route path="/sell" element={<SellPage />} />
        <Route path="/listings/:id" element={<ListingPage />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
);
