import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/context/AuthContext";

const inter = Inter({
  variable: "--font-body",
  subsets: ["latin"],
});

const interDisplay = Inter({
  variable: "--font-display",
  subsets: ["latin"],
  weight: ["600", "700", "800"],
});

const jetbrainsMono = JetBrains_Mono({
  variable: "--font-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "SEAPEDIA — Belanja dari Semua Toko",
  description:
    "SEAPEDIA menghubungkan penjual, pembeli, dan kurir dalam satu marketplace.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id" className={`${inter.variable} ${interDisplay.variable} ${jetbrainsMono.variable} h-full`}>
      <body className="min-h-full flex flex-col bg-sand-50 text-ink antialiased">
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
