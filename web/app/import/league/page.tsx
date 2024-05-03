'use client';

import { useState } from 'react';
import axios from 'axios';

export default function ImportLeaguePage() {
  const [teamsFile, setTeamsFile] = useState<File | null>(null);
  const [playersFile, setPlayersFile] = useState<File | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!teamsFile || !playersFile) {
      alert('Please select both teams and players files.');
      return;
    }

    const formData = new FormData();
    formData.append('teams', teamsFile);
    formData.append('players', playersFile);

    try {
      const response = await axios.post('http://localhost:4000/v1/league/3/import', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      console.log('Import successful:', response.data);
      alert('Import successful!');
    } catch (error) {
      console.error('Import failed:', error);
      alert('Import failed. Please try again.');
    }
  };

  return (
    <div>
      <h1>Import League</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="teams">Teams File:</label>
          <input
            type="file"
            id="teams"
            accept=".csv"
            onChange={(e) => setTeamsFile(e.target.files?.[0] || null)}
          />
        </div>
        <div>
          <label htmlFor="players">Players File:</label>
          <input
            type="file"
            id="players"
            accept=".csv"
            onChange={(e) => setPlayersFile(e.target.files?.[0] || null)}
          />
        </div>
        <button type="submit">Import</button>
      </form>
    </div>
  );
}