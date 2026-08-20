import { NextResponse } from 'next/server';

const GATEWAY = process.env.GATEWAY_URL ?? 'http://localhost:8080';

function extractTokenFromCookie(cookieHeader: string) {
  for (const part of (cookieHeader || '').split(';')) {
    const [name, ...rest] = part.trim().split('=');
    if (name === 'token') {
      return decodeURIComponent(rest.join('=')).trim();
    }
  }
  return '';
}

export async function GET(req: Request) {
  try {
    const cookieHeader = req.headers.get('cookie') || '';
    const token = extractTokenFromCookie(cookieHeader);
    const url = `${GATEWAY}/api/v2/users/me`;
    const res = await fetch(url, {
      method: 'GET',
      headers: {
        Cookie: cookieHeader,
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      cache: 'no-store',
    });

    const body = await res.text();
    return new NextResponse(body, { status: res.status, headers: { 'Content-Type': res.headers.get('Content-Type') || 'application/json' } });
  } catch (err) {
    return NextResponse.json({ error: String(err) }, { status: 500 });
  }
}
