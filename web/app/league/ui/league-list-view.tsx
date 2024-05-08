'use client';

import { League } from "@/app/lib/definitions";
import { redirect, useRouter } from 'next/navigation';
import { redirectUser } from "../actions";

type LeagueListViewProps = {
  leagues: League[] | undefined;
};

const LeagueListView = ({ leagues }: LeagueListViewProps) => {
  const router = useRouter();
  const handleLeagueClick = (leagueId: number) => {
    router.push(`/league/${leagueId}`);
  };

  return (
    <div className="bg-base-100 min-h-screen">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <h1 className="text-5xl font-bold text-neutral-content mb-12">Leagues</h1>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
          {leagues ? (
            leagues.map((league) => (
              <div
                key={league.id}
                onClick={() => handleLeagueClick(league.id)}
                className="bg-base-200 rounded-lg shadow-lg p-8 transform hover:scale-105 transition duration-300 cursor-pointer"
              >
                <div className="flex items-center justify-between mb-6">
                  <div className="text-3xl font-semibold text-neutral-content">{league.name}</div>
                  <div className="text-xl text-error">{league.level}</div>
                </div>
                <div className="border-t border-neutral-content pt-6">
                  <div className="text-lg text-neutral-content mb-2">League Details:</div>
                  <div className="text-base-content">
                    {/* TODO: Add league details here */}
                    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>
                  </div>
                </div>
                <div className="mt-8">
                  <button onClick={() => handleLeagueClick(league.id)} className="btn btn-outline btn-primary btn-block">View League</button>
                </div>
              </div>
            ))
          ) : (
            <div className="text-2xl text-neutral-content">No leagues available</div>
          )}
        </div>
      </div>
    </div>
  );
};

export default LeagueListView;