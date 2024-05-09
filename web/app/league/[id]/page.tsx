import { League } from "@/app/lib/definitions";
import LeagueDetailView from "../ui/league-detail-view";
import { getLeagueById, getTeamsByLeagueId, getPlayersByLeagueId } from "../actions";

export default async function LeagueDetailPage({ params }:{ params: { id: number } }) {
  const league = await getLeagueById(params.id)
  const teams = await getTeamsByLeagueId(params.id)
  const players = await getPlayersByLeagueId(params.id)
  return (
    <LeagueDetailView league={league} teams={teams} players={players} />
  );
}