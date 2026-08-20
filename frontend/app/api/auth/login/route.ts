import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { username, password } = body || {};

    const normalizedUsername = String(username ?? "").trim().toLowerCase();
    const normalizedPassword = String(password ?? "");

    const validCredentials =
      (normalizedUsername === "admin" && normalizedPassword === "password") ||
      (normalizedUsername.includes("@") && normalizedPassword.length >= 8);

    if (validCredentials) {
      const token = `studybuddy-demo-token-${Date.now()}`;
      const displayName = normalizedUsername === "admin" ? "Admin User" : normalizedUsername.split("@")[0];
      const res = NextResponse.json({ success: true, user: { name: displayName, email: normalizedUsername === "admin" ? "admin@studybuddy.local" : normalizedUsername } });

      res.cookies.set({ name: "token", value: token, httpOnly: true, path: "/", maxAge: 60 * 60 * 24 });
      res.cookies.set({ name: "user", value: JSON.stringify({ name: displayName, email: normalizedUsername === "admin" ? "admin@studybuddy.local" : normalizedUsername }), path: "/", maxAge: 60 * 60 * 24 });
      return res;
    }

    return NextResponse.json({ success: false, error: "Invalid credentials" }, { status: 401 });
  } catch (err) {
    return NextResponse.json({ success: false, error: String(err) }, { status: 500 });
  }
}
