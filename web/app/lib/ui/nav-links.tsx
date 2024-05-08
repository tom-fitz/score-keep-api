'use client';

import { UserGroupIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import clsx from 'clsx';

const links = [
  { name: 'Leagues', href: '/league', icon: UserGroupIcon },
];

// const adminLinks = [
//   { name: 'Admin Programs', href: '/admin/program', icon: FolderIcon },
//   { name: 'Admin Exercises', href: '/admin/exercise', icon: RocketLaunchIcon }
// ];

export default function NavLinks() {
  // const { isLoaded, isSignedIn } = useAuth();
  // const { user, isLoaded: userLoaded } = useUser();
  const pathname = usePathname();

  // if (!isLoaded || !isSignedIn || !userLoaded) {
  //   return null;
  // }

  // const isAdmin = user?.publicMetadata?.role === 'admin'; 

  return (
    <>
      {links.map((link) => {
        const LinkIcon = link.icon;
        return (
          <Link
            key={link.name}
            href={link.href}
            className={clsx(
              'flex h-[48px] grow items-center justify-center text-sm font-medium md:flex-none md:justify-start',
              { 'bg-neutral text-base-100': pathname === link.href }
            )}
          >
            <LinkIcon className="w-7 mr-2" />
            {link.name}
          </Link>
        );
      })}
    </>
  );
}