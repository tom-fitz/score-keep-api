'use client';
import { League, Team, Player, PlayerDisplay } from "@/app/lib/definitions";
import { useState } from 'react';
import Papa from 'papaparse';
import { importPlayers, importTeams } from "../actions";
import LeagueEditDialog from "./league-update-dialog";
import { PencilSquareIcon } from "@heroicons/react/24/outline";

type LeagueDetailViewProps = {
  league: League;
  teams?: Team[];
  players?: PlayerDisplay[];
};

export default function LeagueDetailView({ league, teams, players }: LeagueDetailViewProps) {
  const [uploadedTeams, setUploadedTeams] = useState<Team[]>([]);
  const [uploadedPlayers, setUploadedPlayers] = useState<Player[]>([]);
  const [selectedTeam, setSelectedTeam] = useState<Team | null>(null);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);

  console.log("players: ", players)

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

  const handleTeamClick = (team: Team) => {
    setSelectedTeam(team);
  };

  const getPlayersByTeam = (teamId: number) => {
    return players?.filter((player) => player.teamId === teamId);
  };

  const handleEditClick = () => {
    setIsEditDialogOpen(true);
  };

  const handleEditDialogClose = () => {
    setIsEditDialogOpen(false);
  };

  return (
    <div className="bg-gradient-to-br from-base-300 via-base-200 to-base-100 min-h-screen">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="bg-base-200 rounded-lg shadow-lg p-8">
          {/* <div className="flex items-center justify-between mb-8">
            <div className="text-5xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-primary to-accent">{league.name}</div>
            <div className="text-3xl font-bold text-warning">{league.level}</div>
          </div> */}
          <div className="flex items-center justify-between mb-8">
            <div className="text-5xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-primary to-accent">{league.name}</div>
            <div className="flex items-center">
              <div className="text-3xl font-bold text-warning mr-4">{league.level}</div>
              <label htmlFor="edit-league-modal" className="btn btn-ghost btn-circle text-neutral-content hover:text-primary modal-button" onClick={handleEditClick}>
                <PencilSquareIcon className="h-6 w-6" />
              </label>
            </div>
          </div>
          <div className="border-t-4 border-dotted border-neutral-content pt-8">
            <div className="text-3xl font-semibold text-neutral-content mb-4">League Details:</div>
            <div className="text-base-content text-lg">
              {/* TODO: Add detailed league information here */}
              <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus auctor, velit eu tincidunt dictum, odio risus vestibulum tortor, vel sagittis quam tellus quis eros. Sed euismod aliquam neque, id lacinia metus fermentum in. Nulla facilisi. Sed vel erat vel turpis viverra luctus.</p>
            </div>
          </div>

          <div className="mt-16">
            <h2 className="text-4xl font-bold text-neutral-content mb-8">Teams</h2>
            {teams && teams.length > 0 ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-12">
                {teams.map((team) => (
                  <div
                    key={team.id}
                    className={`bg-gradient-to-br from-base-100 to-base-300 rounded-lg shadow-xl p-8 cursor-pointer transition duration-500 ease-in-out transform hover:-translate-y-4 hover:shadow-2xl ${
                      selectedTeam?.id === team.id ? 'ring-4 ring-primary' : ''
                    }`}
                    onClick={() => handleTeamClick(team)}
                  >
                    <div className="text-3xl font-bold text-primary mb-4">{team.name}</div>
                    <div className="text-base-content text-lg mb-2">Captain: {team.captain}</div>
                    <div className="text-base-content text-lg">First Year: {team.firstYear}</div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex items-center justify-center bg-base-100 rounded-lg shadow-xl p-8 h-60">
                <label htmlFor="team-upload" className="btn btn-primary btn-lg">
                  <span className="mr-2 text-lg">Upload Teams</span>
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                  </svg>
                </label>
                <input id="team-upload" type="file" accept=".csv" onChange={handleTeamUpload} className="hidden" />
              </div>
            )}
            {uploadedTeams.length > 0 && (
              <div className="mt-12">
                <h3 className="text-3xl font-semibold text-neutral-content mb-6">Uploaded Teams</h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-12">
                  {uploadedTeams.map((team) => (
                    <div
                      key={team.id}
                      className={`bg-gradient-to-br from-base-100 to-base-300 rounded-lg shadow-xl p-8 cursor-pointer transition duration-500 ease-in-out transform hover:-translate-y-4 hover:shadow-2xl ${
                        selectedTeam?.id === team.id ? 'ring-4 ring-primary' : ''
                      }`}
                      onClick={() => handleTeamClick(team)}
                    >
                      <div className="text-3xl font-bold text-primary mb-4">{team.name}</div>
                      <div className="text-base-content text-lg mb-2">Captain: {team.captain}</div>
                      <div className="text-base-content text-lg">First Year: {team.firstYear}</div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          {selectedTeam && (
            <div className="mt-16">
              <h2 className="text-4xl font-bold text-neutral-content mb-8">Players - {selectedTeam.name}</h2>
              {getPlayersByTeam(selectedTeam.id)?.length ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-12">
                  {getPlayersByTeam(selectedTeam.id)?.map((player) => (
                    <div key={player.id} className="bg-base-100 rounded-lg shadow-xl overflow-hidden">
                      <div className="bg-gradient-to-r from-primary to-secondary h-40 flex items-center justify-center">
                        <div className="text-4xl font-extrabold text-base-100">
                          {player.firstName.charAt(0)}
                          {player.lastName.charAt(0)}
                        </div>
                      </div>
                      <div className="p-6">
                        <div className="text-2xl font-bold text-neutral-content mb-2">
                          {player.firstName} {player.lastName}
                        </div>
                        <div className="text-base-content text-lg mb-1">Email: {player.email}</div>
                        <div className="text-base-content text-lg mb-1">Phone: {player.phone}</div>
                        <div className="text-base-content text-lg mb-1">Level: {player.level}</div>
                        <div className="text-base-content text-lg">USA#: {player.usanum}</div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-base-content text-lg">No players available for this team.</div>
              )}
            </div>
          )}

          <div className="mt-16">
            <h2 className="text-4xl font-bold text-neutral-content mb-8">All Players</h2>
            {players && players.length > 0 ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-12">
                {players.map((player) => (
                  <div key={player.id} className="bg-base-100 rounded-lg shadow-xl overflow-hidden">
                    <div className="bg-gradient-to-r from-accent to-warning h-40 flex items-center justify-center">
                      <div className="text-4xl font-extrabold text-base-100">
                        {player.firstName.charAt(0)}
                        {player.lastName.charAt(0)}
                      </div>
                    </div>
                    <div className="p-6">
                      <div className="text-2xl font-bold text-neutral-content mb-2">
                        {player.firstName} {player.lastName}
                      </div>
                      <div className="text-base-content text-lg mb-1">Email: {player.email}</div>
                      <div className="text-base-content text-lg mb-1">Phone: {player.phone}</div>
                      <div className="text-base-content text-lg mb-1">Level: {player.level}</div>
                      <div className="text-base-content text-lg mb-1">USA#: {player.usanum}</div>
                      <div className="text-base-content text-lg">Team: {player.teamName}</div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex items-center justify-center bg-base-100 rounded-lg shadow-xl p-8 h-60">
                <label htmlFor="player-upload" className="btn btn-primary btn-lg">
                  <span className="mr-2 text-lg">Upload Players</span>
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                  </svg>
                </label>
                <input id="player-upload" type="file" accept=".csv" onChange={handlePlayerUpload} className="hidden" />
              </div>
            )}
            {uploadedPlayers.length > 0 && (
              <div className="mt-12">
                <h3 className="text-3xl font-semibold text-neutral-content mb-6">Uploaded Players</h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-12">
                  {uploadedPlayers.map((player) => (
                    <div key={player.id} className="bg-base-100 rounded-lg shadow-xl overflow-hidden">
                      <div className="bg-gradient-to-r from-secondary to-accent h-40 flex items-center justify-center">
                        <div className="text-4xl font-extrabold text-base-100">
                          {player.firstName.charAt(0)}
                          {player.lastName.charAt(0)}
                        </div>
                      </div>
                      <div className="p-6">
                        <div className="text-2xl font-bold text-neutral-content mb-2">
                          {player.firstName} {player.lastName}
                        </div>
                        <div className="text-base-content text-lg mb-1">Email: {player.email}</div>
                        <div className="text-base-content text-lg mb-1">Phone: {player.phone}</div>
                        <div className="text-base-content text-lg mb-1">Level: {player.level}</div>
                        <div className="text-base-content text-lg">USA#: {player.usanum}</div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
      <LeagueEditDialog
        isOpen={isEditDialogOpen}
        onClose={handleEditDialogClose}
        league={league}
      />
    </div>
  );
}