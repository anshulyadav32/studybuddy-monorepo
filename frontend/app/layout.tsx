import './globals.css';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'StudyBuddy Student App',
  description: 'Student learning dashboard',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
