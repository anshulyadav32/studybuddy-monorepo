import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const provider = String(body?.provider ?? "google");
    const action = String(body?.action ?? "signin");

    if (provider !== "google") {
      return NextResponse.json({ success: false, error: "Unsupported auth provider." }, { status: 400 });
    }

    const displayName = action === "signup" ? "Google User" : "Google User";
    const email = "google-user@studybuddy.local";
    const token = `studybuddy-google-token-${Date.now()}`;
    const res = NextResponse.json({ success: true, user: { name: displayName, email }, provider });

    res.cookies.set({ name: "token", value: token, httpOnly: true, path: "/", maxAge: 60 * 60 * 24 });
    res.cookies.set({ name: "user", value: JSON.stringify({ name: displayName, email }), path: "/", maxAge: 60 * 60 * 24 });
    return res;
  } catch (err) {
    return NextResponse.json({ success: false, error: String(err) }, { status: 500 });
  }
}
