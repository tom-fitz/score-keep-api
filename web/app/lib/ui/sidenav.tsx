'use client'

import NavLinks from '@/app/lib/ui/nav-links';
// import { UserButton } from '@clerk/nextjs';


export default function SideNav() {
  return (
    <div className="flex h-full flex-col bg-base-300 p-4">
      <div className="p-2 flex grow flex-row md:flex-col md:space-x-0 md:space-y-2">
        <div className="ma-2">
          <NavLinks />
        </div>
        <div className="hidden h-auto w-full grow rounded-md bg-natural md:block"></div>
        {/* <div className='ma-2'>
          <UserButton />
        </div> */}
      </div>
    </div>
  );
}
