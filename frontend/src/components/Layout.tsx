import type { ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const isDetailPage = location.pathname.startsWith('/funds/');

  return (
    <div className="min-h-screen bg-dark-950">
      {/* Background gradient effects */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-primary-500/10 rounded-full blur-3xl" />
        <div className="absolute top-1/2 -left-40 w-80 h-80 bg-accent-500/10 rounded-full blur-3xl" />
        <div className="absolute -bottom-40 right-1/3 w-80 h-80 bg-primary-500/5 rounded-full blur-3xl" />
      </div>

      {/* Skip link for accessibility */}
      <a
        href="#main-content"
        className="skip-link"
      >
        Skip to main content
      </a>

      {/* Header */}
      <header className="sticky top-0 z-40 backdrop-blur-xl bg-dark-950/80 border-b border-white/[0.06]">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <nav className="flex items-center justify-between h-16">
            <Link
              to="/"
              className="group flex items-center gap-3 focus-ring rounded-xl px-2 py-1 -ml-2"
            >
              {/* Logo icon */}
              <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center shadow-glow-sm group-hover:shadow-glow transition-shadow duration-300">
                <svg className="w-5 h-5 text-dark-900" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <span className="text-lg font-heading font-semibold text-slate-100 group-hover:text-primary-400 transition-colors">
                  Augment
                </span>
                <span className="text-lg font-heading font-semibold text-primary-500">
                  {' '}Fund
                </span>
              </div>
            </Link>

            <div className="flex items-center gap-4">
              {isDetailPage && (
                <Link
                  to="/"
                  className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 focus-ring rounded-lg px-3 py-2 transition-colors"
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
        className="relative mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8"
        tabIndex={-1}
      >
        {children}
      </main>
    </div>
  );
}
