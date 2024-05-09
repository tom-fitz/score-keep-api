'use client';
import { League, Team, Player } from "@/app/lib/definitions";
import { useState } from 'react';
import Papa from 'papaparse';
import { importPlayers, importTeams } from "../actions";

type LeagueDetailViewProps = {
  league: League;
  teams?: Team[];
  players?: Player[];
};

export default function LeagueDetailView({ league, teams, players }: LeagueDetailViewProps) {
  const [uploadedTeams, setUploadedTeams] = useState<Team[]>([]);
  const [uploadedPlayers, setUploadedPlayers] = useState<Player[]>([]);

  const handleTeamUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      Papa.parse(file, {
        header: true,
        complete: (results) => {
          const newTeams: Team[] = results.data.map((row: any, index) => ({
            id: index + 1,
            name: row.name,
            captain: row.captain,
            firstYear: row.firstYear,
          }));
          setUploadedTeams(newTeams);
          importTeams(file)
        },
      });
    }
  };

  const handlePlayerUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      Papa.parse(file, {
        header: true,
        complete: (results) => {
          const newPlayers: Player[] = results.data.map((row: any, index) => ({
            id: index + 1,
            firstName: row.firstName,
            lastName: row.lastName,
            email: row.email,
            phone: row.phone,
            usanum: row.usaNum,
            level: row.level,
          }));
          setUploadedPlayers(newPlayers);
          importPlayers(file)
        },
      });
    }
  };

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

          <div className="mt-8">
            <h2 className="text-2xl font-semibold text-neutral-content mb-4">Teams</h2>
            {teams && teams.length > 0 ? (
              <ul className="space-y-2">
                {teams.map((team) => (
                  <li key={team.id} className="text-base-content">
                    {team.name}
                  </li>
                ))}
              </ul>
            ) : (
              <div>
                <input placeholder="No teams available. Upload teams" type="file" accept=".csv" onChange={handleTeamUpload} />
              </div>
            )}
            {uploadedTeams.length > 0 && (
              <div className="mt-4">
                <h3 className="text-xl font-semibold text-neutral-content mb-2">Uploaded Teams</h3>
                <ul className="space-y-2">
                  {uploadedTeams.map((team) => (
                    <li key={team.id} className="text-base-content">
                      {team.name}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>

          <div className="mt-8">
            <h2 className="text-2xl font-semibold text-neutral-content mb-4">Players</h2>
            {players && players.length > 0 ? (
              <ul className="space-y-2">
                {players.map((player) => (
                  <li key={player.id} className="text-base-content">
                    {player.firstName} {player.lastName}
                  </li>
                ))}
              </ul>
            ) : (
              <div>
                <input placeholder="No players available. Upload players" type="file" accept=".csv" onChange={handlePlayerUpload} />
              </div>
            )}
            {uploadedPlayers.length > 0 && (
              <div className="mt-4">
                <h3 className="text-xl font-semibold text-neutral-content mb-2">Uploaded Players</h3>
                <ul className="space-y-2">
                  {uploadedPlayers.map((player) => (
                    <li key={player.id} className="text-base-content">
                      {player.firstName} {player.lastName}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}