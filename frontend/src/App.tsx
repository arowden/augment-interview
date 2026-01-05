import { Routes, Route } from 'react-router-dom';

import { Layout } from './components/Layout';
import { ErrorBoundary } from './components/ErrorBoundary';
import { Dashboard } from './pages/Dashboard';
import { FundPage } from './pages/FundPage';
import { OwnersPage } from './pages/OwnersPage';
import { NotFound } from './pages/NotFound';

function App() {
  return (
    <Layout>
      <ErrorBoundary>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/funds/:id" element={<FundPage />} />
          <Route path="/owners" element={<OwnersPage />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </ErrorBoundary>
    </Layout>
  );
}

export default App;
