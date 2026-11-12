import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import type { User } from "../types";

export default function ProfilePage() {
  const [user, setUser] = useState<User | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      navigate("/login");
      return;
    }

    fetch("/api/me", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (!res.ok) {
          navigate("/login");
          return null;
        }
        return res.json();
      })
      .then((data) => {
        if (data) {
          setUser(data);
        }
      });
  }, [navigate]);

  if (user === null) {
    return (
      <main>
        <p>Loading...</p>
      </main>
    );
  }

  return (
    <main className="profile-page">
      <div className="avatar profile-avatar">{user.name[0].toUpperCase()}</div>
      <h1>{user.name}</h1>
      <p className="profile-email">{user.email}</p>
    </main>
  );
}
