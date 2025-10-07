"use client";

import { useState } from "react";
import { usePrivy } from "@privy-io/react-auth";
import { Web3Assets } from "../components/Web3Assets";

// Define types for the analysis response based on the Swagger documentation
interface AnalysisResponse {
  analysis: {
    security: {
      total_issues: number;
      findings: any[]; // Replace 'any' with a specific type if available
    };
  };
  suggestions: any[]; // Replace 'any' with a specific type if available
  score: number;
}

export default function Home() {
  const { ready, authenticated, login, logout } = usePrivy();
  const [code, setCode] = useState("");
  const [analysis, setAnalysis] = useState<AnalysisResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAnalysis = async () => {
    if (!authenticated) {
      setError("Please log in to analyze your code.");
      return;
    }
    setIsLoading(true);
    setError(null);
    setAnalysis(null);
    try {
      const response = await fetch("/api/analyze", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content: code }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to analyze the code.");
      }

      const data: AnalysisResponse = await response.json();
      setAnalysis(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An unknown error occurred.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-gray-900 text-white min-h-screen flex flex-col">
      <header className="fixed top-0 left-0 right-0 bg-gray-900 bg-opacity-50 backdrop-blur-md z-10">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex-shrink-0">
              <h1 className="text-xl font-bold">IaC AI Agent</h1>
            </div>
            <div className="flex items-center">
              {ready && (
                <button
                  onClick={authenticated ? logout : login}
                  className="bg-blue-600 px-4 py-2 text-white font-semibold rounded-lg shadow-md hover:bg-blue-700 transition duration-300"
                >
                  {authenticated ? "Logout" : "Login"}
                </button>
              )}
            </div>
          </div>
        </div>
      </header>

      <main className="flex-grow container mx-auto px-4 sm:px-6 lg:px-8 pt-24 pb-12">
        <section className="text-center">
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-blue-600">
            Secure Your Infrastructure with AI
          </h1>
          <p className="mt-4 max-w-2xl mx-auto text-lg text-gray-400">
            Paste your Terraform code, and our AI-powered agent will analyze it for security vulnerabilities, cost optimizations, and best practices in seconds.
          </p>
        </section>

        {authenticated && (
          <section className="mt-8">
            <Web3Assets />
          </section>
        )}

        <section className="mt-8 bg-gray-800 p-6 rounded-xl shadow-lg">
          <h2 className="text-2xl font-bold mb-4">Terraform Code Analysis</h2>
          <textarea
            className="w-full h-64 p-4 font-mono text-sm bg-gray-900 border border-gray-700 rounded-md text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-600 transition"
            placeholder={authenticated ? "Paste your Terraform code here..." : "Please log in to use the analysis tool."}
            value={code}
            onChange={(e) => setCode(e.target.value)}
            disabled={!authenticated}
          />
          <div className="mt-4 flex justify-end">
            <button
              className="px-8 py-3 bg-blue-600 text-white font-semibold rounded-lg shadow-md hover:bg-blue-700 disabled:bg-gray-500 disabled:cursor-not-allowed transition duration-300"
              onClick={handleAnalysis}
              disabled={isLoading || !code || !authenticated}
            >
              {isLoading ? "Analyzing..." : "Analyze Code"}
            </button>
          </div>
        </section>

        {error && (
          <section className="mt-6 text-center">
            <p className="text-red-400 bg-red-900/50 p-4 rounded-lg">{error}</p>
          </section>
        )}

        {isLoading && (
            <div className="text-center mt-6">
                <p className="text-lg">Analyzing your code, please wait...</p>
            </div>
        )}

        {analysis && (
          <section className="mt-6 bg-gray-800 p-6 rounded-xl shadow-lg">
            <h3 className="text-2xl font-bold mb-4">Analysis Results</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="bg-gray-700 p-4 rounded-lg">
                <h4 className="text-xl font-semibold">Overall Score</h4>
                <p className="text-4xl font-bold text-green-400">{analysis.score}/100</p>
              </div>
              <div className="bg-gray-700 p-4 rounded-lg">
                <h4 className="text-xl font-semibold">Security Issues</h4>
                <p className="text-4xl font-bold text-yellow-400">{analysis.analysis.security.total_issues}</p>
              </div>
            </div>
            <div className="mt-6">
              <h4 className="text-xl font-semibold">Suggestions</h4>
              <ul className="mt-2 space-y-2">
                {analysis.suggestions.map((suggestion, index) => (
                  <li key={index} className="p-3 bg-gray-700 rounded-lg">
                    {suggestion.message}
                  </li>
                ))}
              </ul>
            </div>
          </section>
        )}
      </main>

      <footer className="text-center py-6 text-gray-500">
        <p>Powered by the IaC AI Agent</p>
      </footer>
    </div>
  );
}