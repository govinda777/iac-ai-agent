"use client";

import { useAccount, useBalance, useReadContract } from "wagmi";
import { formatUnits } from "viem";

// IMPORTANT: Replace these with your actual contract addresses and ABIs
const NFT_CONTRACT_ADDRESS = "0x...YOUR_NFT_CONTRACT_ADDRESS";
const TOKEN_CONTRACT_ADDRESS = "0x...YOUR_TOKEN_CONTRACT_ADDRESS";

const NFT_ABI = [
  {
    type: "function",
    name: "balanceOf",
    stateMutability: "view",
    inputs: [{ name: "owner", type: "address" }],
    outputs: [{ type: "uint256" }],
  },
] as const;

export function Web3Assets() {
  const { address, isConnected } = useAccount();

  const { data: nftBalance, isLoading: isNftLoading } = useReadContract({
    address: NFT_CONTRACT_ADDRESS,
    abi: NFT_ABI,
    functionName: "balanceOf",
    args: [address!],
    query: {
      enabled: isConnected && !!address,
    },
  });

  const { data: tokenBalance, isLoading: isTokenLoading } = useBalance({
    address,
    token: TOKEN_CONTRACT_ADDRESS,
    query: {
      enabled: isConnected && !!address,
    },
  });

  const getAccessTier = (balance: bigint | undefined) => {
    if (!isConnected || balance === undefined) return "Loading...";
    if (balance === 0n) return "No Access NFT";
    if (balance > 0n) return "Pro Access"; // Example tier
    return "Unknown";
  };

  if (!isConnected) {
    return (
      <div className="p-4 bg-gray-800 border border-gray-700 rounded-md text-center">
        <p>Connect your wallet to see your assets.</p>
      </div>
    );
  }

  return (
    <div className="p-4 bg-gray-800 border border-gray-700 rounded-md">
      <h3 className="text-xl font-bold mb-4 text-center">Your Web3 Assets</h3>
      <div className="space-y-2">
        <div>
          <strong>Access Tier:</strong>{" "}
          <span>{isNftLoading ? "Loading..." : getAccessTier(nftBalance)}</span>
        </div>
        <div>
          <strong>IACAI Balance:</strong>{" "}
          <span>
            {isTokenLoading
              ? "Loading..."
              : tokenBalance
              ? `${formatUnits(tokenBalance.value, tokenBalance.decimals)} ${
                  tokenBalance.symbol
                }`
              : "0"}
          </span>
        </div>
      </div>
    </div>
  );
}