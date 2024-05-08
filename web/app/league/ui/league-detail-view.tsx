'use client';

import { League } from "@/app/lib/definitions";

export default function LeagueDetailView({ league }:{ league: League}){
    console.log(league);
    return (
        <div className="bg-base-100 min-h-screen">
            <div className="max-w-7xl mx-auto px-4 py-8">
                <div className="bg-base-200 rounded-lg shadow-lg p-8">
                <div className="flex items-center justify-between mb-6">
                    <div className="text-3xl font-semibold text-neutral-content">{league.name}</div>
                    <div className="text-xl text-error">{league.level}</div>
                </div>
                <div className="border-t border-neutral-content pt-6">
                    <div className="text-lg text-neutral-content mb-2">League Details:</div>
                    <div className="text-base-content">
                    {/* TODO: Add detailed league information here */}
                    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>
                    </div>
                </div>
                {/* TODO: Add more sections or components specific to the league detail view */}
                </div>
            </div>
        </div>
    )
}