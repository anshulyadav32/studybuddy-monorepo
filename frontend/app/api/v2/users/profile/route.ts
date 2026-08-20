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

async function proxy(req: Request) {
  const cookieHeader = req.headers.get('cookie') || '';
  const token = extractTokenFromCookie(cookieHeader);
  const url = `${GATEWAY}/api/v2/users/profile`;

  const init: RequestInit = {
    method: req.method,
    headers: {
      Cookie: cookieHeader,
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      'Content-Type': req.headers.get('content-type') || 'application/json',
    },
    body: ['GET', 'HEAD'].includes(req.method) ? undefined : await req.text(),
    redirect: 'manual',
  };

  const res = await fetch(url, init);
  const body = await res.text();
  return new NextResponse(body, { status: res.status, headers: { 'Content-Type': res.headers.get('Content-Type') || 'application/json' } });
}

export async function GET(req: Request) {
  try {
    return await proxy(req);
  } catch (err) {
    return NextResponse.json({ error: String(err) }, { status: 500 });
  }
}

export async function PUT(req: Request) {
  try {
    return await proxy(req);
  } catch (err) {
    return NextResponse.json({ error: String(err) }, { status: 500 });
  }
}

export async function PATCH(req: Request) {
  try {
    return await proxy(req);
  } catch (err) {
    return NextResponse.json({ error: String(err) }, { status: 500 });
  }
}
