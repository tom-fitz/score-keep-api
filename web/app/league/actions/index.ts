import { League, Player, Team } from '@/app/lib/definitions';
import axios from 'axios';
import { redirect } from 'next/navigation';

export function getLeagues(): Promise<League[]>{
    try {
        return axios.get('http://localhost:4000/v1/league').then((response) => response.data.leagues);
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

export function importTeams(file: File): Promise<void> {
    try {
        const formData = new FormData();
        formData.append('teams', file, 'teams.csv');

        return fetch('http://localhost:4000/v1/league/1/teams/import', {
            method: 'POST',
            body: formData,
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to import teams');
            }
            return response.json();
        })
        .then(result => {
            console.log(result);
        })
        .catch(error => {
            console.error('Error importing teams:', error);
            throw error;
        });
    } catch (error) {
        throw error;
    }
}

export function importPlayers(file: File): Promise<void> {
    try {
        const formData = new FormData();
        formData.append('players', file, 'players.csv');

        return fetch('http://localhost:4000/v1/league/1/players/import', {
            method: 'POST',
            body: formData,
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to import players');
            }
            return response.json();
        })
        .then(result => {
            console.log(result);
        })
        .catch(error => {
            console.error('Error importing players:', error);
            throw error;
        });
    } catch (error) {
        throw error;
    }
}

export function getTeamsByLeagueId(lid: number): Promise<Team[]>{
    try {
        return axios.get(`http://localhost:4000/v1/league/${lid}/teams`).then((response) => response.data.teams)
    } catch (error) {
        throw error;
    }
}

export function getPlayersByLeagueId(lid: number): Promise<Player[]>{
    try {
        return axios.get(`http://localhost:4000/v1/league/${lid}/players`).then((response) => response.data.players)
    } catch (error) {
        throw error;
    }
}

