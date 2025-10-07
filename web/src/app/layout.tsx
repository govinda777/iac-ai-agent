import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Providers } from "./providers";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "IaC AI Agent",
  description: "AI-powered assistant for Terraform code analysis",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body
        className={`${inter.className} bg-gray-900 text-gray-100 flex flex-col items-center justify-center min-h-screen`}
      >
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}