import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const name = String(body?.name ?? "").trim();
    const email = String(body?.email ?? "").trim();
    const password = String(body?.password ?? "");

    if (!name || !email || !password) {
      return NextResponse.json({ success: false, error: "Name, email, and password are required." }, { status: 400 });
    }

    if (password.length < 8) {
      return NextResponse.json({ success: false, error: "Password must be at least 8 characters long." }, { status: 400 });
    }

    const token = `studybuddy-demo-token-${Date.now()}`;
    const res = NextResponse.json({ success: true, user: { name, email } });

    res.cookies.set({ name: "token", value: token, httpOnly: true, path: "/", maxAge: 60 * 60 * 24 });
    res.cookies.set({ name: "user", value: JSON.stringify({ name, email }), path: "/", maxAge: 60 * 60 * 24 });
    return res;
  } catch (err) {
    return NextResponse.json({ success: false, error: String(err) }, { status: 500 });
  }
}
