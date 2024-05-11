'use client'

import { League } from '@/app/lib/definitions';
import { updateLeague } from '../actions';

type LeagueEditDialogProps = {
  isOpen: boolean;
  onClose: () => void;
  league: League;
};

export default function LeagueEditDialog({ isOpen, onClose, league }: LeagueEditDialogProps) {
  const handleSave = () => {
    // Handle saving the edited league details
    // updateLeague(league)
    console.log(league)
    // onClose();
  };

  return (
    <>
      <input type="checkbox" id="edit-league-modal" className="modal-toggle" checked={isOpen} readOnly />
      <div className="modal">
        <div className="modal-box">
          <h3 className="font-bold text-2xl mb-4">Edit League</h3>
          <div className="mb-4">
            <label htmlFor="name" className="block mb-2">Name</label>
            <input
              type="text"
              name="name"
              id="name"
              className="input input-bordered w-full"
              defaultValue={league.name}
            />
          </div>
          <div className="mb-4">
            <label htmlFor="level" className="block mb-2">Level</label>
            <input
              type="text"
              name="level"
              id="level"
              className="input input-bordered w-full"
              defaultValue={league.level}
            />
          </div>
          <div className="modal-action">
            <label htmlFor="edit-league-modal" className="btn btn-primary" onClick={handleSave}>Save</label>
            <label htmlFor="edit-league-modal" className="btn" onClick={onClose}>Cancel</label>
          </div>
        </div>
      </div>
    </>
  );
}