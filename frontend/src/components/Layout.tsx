import type { ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const isDetailPage = location.pathname.startsWith('/funds/');
  const isOwnersPage = location.pathname === '/owners';

  return (
    <div className="min-h-screen bg-slate-950">
      {/* Background gradient effects */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-96 h-96 bg-indigo-500/5 rounded-full blur-3xl" />
        <div className="absolute top-1/3 -left-40 w-96 h-96 bg-violet-500/4 rounded-full blur-3xl" />
        <div className="absolute -bottom-40 right-1/4 w-96 h-96 bg-purple-500/3 rounded-full blur-3xl" />
      </div>

      {/* Skip link for accessibility */}
      <a href="#main-content" className="skip-link">
        Skip to main content
      </a>

      {/* Header */}
      <header className="sticky top-0 z-40 backdrop-blur-xl bg-slate-950/80 border-b border-slate-800/50">
        <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
          <nav className="flex items-center justify-between h-16">
            <Link
              to="/"
              className="group flex items-center gap-3 focus-ring rounded-xl px-2 py-1 -ml-2"
            >
              <div className="w-8 h-8 rounded-lg bg-indigo-500/20 flex items-center justify-center">
                <svg className="w-5 h-5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 3v11.25A2.25 2.25 0 006 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0118 16.5h-2.25m-7.5 0h7.5m-7.5 0l-1 3m8.5-3l1 3m0 0l.5 1.5m-.5-1.5h-9.5m0 0l-.5 1.5M9 11.25v1.5M12 9v3.75m3-6v6" />
                </svg>
              </div>
              <span className="text-xl font-semibold tracking-tight text-white">
                Cap Table
              </span>
            </Link>

            <div className="flex items-center gap-3">
              <Link
                to="/owners"
                className={`flex items-center gap-2 text-sm rounded-lg px-3 py-2 transition-all ${
                  isOwnersPage
                    ? 'text-white bg-indigo-600'
                    : 'text-slate-400 hover:text-white bg-slate-800/50 hover:bg-slate-800'
                }`}
              >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
                </svg>
                Owners
              </Link>
              {(isDetailPage || isOwnersPage) && (
                <Link
                  to="/"
                  className="flex items-center gap-2 text-sm text-slate-400 hover:text-white bg-slate-800/50 hover:bg-slate-800 rounded-lg px-3 py-2 transition-all"
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                  </svg>
                  Dashboard
                </Link>
              )}
            </div>
          </nav>
        </div>
      </header>

      {/* Main content */}
      <main
        id="main-content"
        className="relative mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8"
        tabIndex={-1}
      >
        {children}
      </main>
    </div>
  );
}
