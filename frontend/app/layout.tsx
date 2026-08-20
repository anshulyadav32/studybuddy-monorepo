import './globals.css';
import type { Metadata } from 'next';
import Header from './components/Header';
import { cookies } from 'next/headers';

export const metadata: Metadata = {
  title: 'StudyBuddy Student App',
  description: 'Student learning dashboard',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  // Read "user" cookie on the server and pass to the client Header component
  const cookieStore = cookies();
  const userCookie = cookieStore.get('user')?.value;
  let user: { name?: string; email?: string } | null = null;
  try {
    user = userCookie ? JSON.parse(userCookie) : null;
  } catch (e) {
    user = null;
  }

  return (
    <html lang="en">
      <body>
        <Header user={user} />
        {children}
      </body>
    </html>
  );
}
