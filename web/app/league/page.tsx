import { getLeagues } from "./actions";
import LeagueListView from "./ui/league-list-view";

export default async function LeagueLayout(){
    const leagues = await getLeagues();
    console.log(leagues);
    return (
        <div className={'container align-center justify-center text-center h-full w-full'}>
            <LeagueListView leagues={leagues} />
        </div>
    )
}