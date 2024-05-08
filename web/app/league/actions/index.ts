import { League } from '@/app/lib/definitions';
import axios from 'axios';
import { redirect } from 'next/navigation';

export function getLeagues(): Promise<League[]>{
    try {
        const api = 'http://localhost:4000/v1/league';
        return axios.get(api).then((response) => response.data.leagues);
    } catch(error) {
        throw error;
    }
}

export function getLeagueById(id: number): Promise<League>{
    try {
        return axios.get(`http://localhost:4000/v1/league/${id}`).then((response) => response.data.league)
    } catch (error) {
        throw error;
    }
}

export function redirectUser(url: string) {
    try {
        return redirect(url)
    } catch (error) {
        throw error;
    }
}