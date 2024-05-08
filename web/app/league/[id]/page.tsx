import { League } from "@/app/lib/definitions";
import LeagueDetailView from "../ui/league-detail-view";
import { getLeagueById } from "../actions";

export default async function LeagueDetailPage({ params }:{ params: { id: number } }) {
  const league = await getLeagueById(params.id)
  return (
    <LeagueDetailView league={league} />
  );
}