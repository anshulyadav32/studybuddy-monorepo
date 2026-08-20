import { NextResponse } from "next/server";

export async function POST() {
  try {
    const res = NextResponse.json({ success: true });

    // Clear cookies by setting empty value and maxAge 0 / expires
    res.cookies.set({ name: "token", value: "", path: "/", maxAge: 0 });
    res.cookies.set({ name: "user", value: "", path: "/", maxAge: 0 });

    return res;
  } catch (err) {
    return NextResponse.json({ success: false, error: String(err) }, { status: 500 });
  }
}
