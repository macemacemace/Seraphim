import { useEffect, useState } from "react";
import type { User } from "./types";

export function useCurrentUser() {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      return;
    }

    fetch("/api/me", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (!res.ok) {
          return null;
        }
        return res.json();
      })
      .then((data) => {
        if (data) {
          setUser(data);
        }
      });
  }, []);

  return user;
}
