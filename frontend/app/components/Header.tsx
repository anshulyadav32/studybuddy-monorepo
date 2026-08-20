"use client";

import React from "react";

type User = { name?: string; email?: string } | null;

export default function Header({ user }: { user: User }) {
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  async function signOut() {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch("/api/auth/logout", { method: "POST" });
      const data = await res.json();
      if (res.ok && data.success) {
        // redirect to login page
        window.location.href = "/login";
      } else {
        setError(data?.error || "Sign out failed");
      }
    } catch (err) {
      setError(String(err));
    } finally {
      setLoading(false);
    }
  }

  return (
    <header style={{ background: "#ffffff", borderBottom: "1px solid #e6edf3", padding: "12px 20px" }}>
      <div style={{ maxWidth: 1200, margin: "0 auto", display: "flex", alignItems: "center", justifyContent: "space-between" }}>
        <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
          <div style={{ fontWeight: 800, color: "#2563eb" }}>StudyBuddy</div>
        </div>

        <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
          {user ? (
            <>
              <div style={{ textAlign: "right", lineHeight: 1 }}>
                <div style={{ fontWeight: 700 }}>{user.name || user.email}</div>
                <div style={{ fontSize: 12, color: "#64748b" }}>{user.email}</div>
              </div>
              <button
                onClick={signOut}
                disabled={loading}
                style={{ background: "#ef4444", color: "#fff", border: "none", padding: "8px 12px", borderRadius: 8, fontWeight: 700 }}
              >
                {loading ? "Signing out…" : "Sign out"}
              </button>
            </>
          ) : (
            <a href="/login" style={{ fontWeight: 700, color: "#2563eb" }}>Sign in</a>
          )}
        </div>
      </div>
      {error && <div style={{ maxWidth: 1200, margin: "6px auto 0", color: "#b91c1c", fontSize: 13 }}>{error}</div>}
    </header>
  );
}
