import { http, HttpResponse } from 'msw';

import type { Fund, CapTable, Transfer } from '../api/client';

export const mockFunds: Fund[] = [
  {
    id: '550e8400-e29b-41d4-a716-446655440000',
    name: 'Growth Fund I',
    totalUnits: 1000000,
    createdAt: '2024-01-15T10:30:00Z',
  },
  {
    id: '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
    name: 'Venture Fund II',
    totalUnits: 500000,
    createdAt: '2024-02-20T14:45:00Z',
  },
];

export const mockCapTable: CapTable = {
  fundId: '550e8400-e29b-41d4-a716-446655440000',
  entries: [
    {
      ownerName: 'Founder LLC',
      units: 600000,
      percentage: 60.0,
      acquiredAt: '2024-01-15T10:30:00Z',
    },
    {
      ownerName: 'Investor A',
      units: 250000,
      percentage: 25.0,
      acquiredAt: '2024-03-01T09:00:00Z',
    },
    {
      ownerName: 'Investor B',
      units: 150000,
      percentage: 15.0,
      acquiredAt: '2024-03-15T11:30:00Z',
    },
  ],
  total: 3,
  limit: 100,
  offset: 0,
};

export const mockTransfers: Transfer[] = [
  {
    id: '7c9e6679-7425-40de-944b-e07fc1f90ae7',
    fundId: '550e8400-e29b-41d4-a716-446655440000',
    fromOwner: 'Founder LLC',
    toOwner: 'Investor A',
    units: 250000,
    transferredAt: '2024-03-01T09:00:00Z',
  },
  {
    id: '8d0e6680-8526-51ef-a55c-f18fd2g01bf8',
    fundId: '550e8400-e29b-41d4-a716-446655440000',
    fromOwner: 'Founder LLC',
    toOwner: 'Investor B',
    units: 150000,
    transferredAt: '2024-03-15T11:30:00Z',
  },
];

export const handlers = [
  http.get('/api/funds', () => {
    return HttpResponse.json(mockFunds);
  }),

  http.post('/api/funds', async ({ request }) => {
    const body = await request.json() as { name: string; totalUnits: number };
    const newFund: Fund = {
      id: crypto.randomUUID(),
      name: body.name,
      totalUnits: body.totalUnits,
      createdAt: new Date().toISOString(),
    };
    return HttpResponse.json(newFund, { status: 201 });
  }),

  http.get('/api/funds/:fundId', ({ params }) => {
    const fund = mockFunds.find((f) => f.id === params['fundId']);
    if (!fund) {
      return HttpResponse.json(
        { code: 'FUND_NOT_FOUND', message: 'Fund not found' },
        { status: 404 }
      );
    }
    return HttpResponse.json(fund);
  }),

  http.get('/api/funds/:fundId/cap-table', ({ params }) => {
    const fund = mockFunds.find((f) => f.id === params['fundId']);
    if (!fund) {
      return HttpResponse.json(
        { code: 'FUND_NOT_FOUND', message: 'Fund not found' },
        { status: 404 }
      );
    }
    return HttpResponse.json({ ...mockCapTable, fundId: params['fundId'] });
  }),

  http.get('/api/funds/:fundId/transfers', ({ params }) => {
    const fund = mockFunds.find((f) => f.id === params['fundId']);
    if (!fund) {
      return HttpResponse.json(
        { code: 'FUND_NOT_FOUND', message: 'Fund not found' },
        { status: 404 }
      );
    }
    return HttpResponse.json(
      mockTransfers.filter((t) => t.fundId === params['fundId'])
    );
  }),

  http.post('/api/funds/:fundId/transfers', async ({ params, request }) => {
    const fund = mockFunds.find((f) => f.id === params['fundId']);
    if (!fund) {
      return HttpResponse.json(
        { code: 'FUND_NOT_FOUND', message: 'Fund not found' },
        { status: 404 }
      );
    }

    const body = await request.json() as { fromOwner: string; toOwner: string; units: number };

    if (body.fromOwner === body.toOwner) {
      return HttpResponse.json(
        { code: 'SELF_TRANSFER', message: 'Cannot transfer units to yourself' },
        { status: 400 }
      );
    }

    const newTransfer: Transfer = {
      id: crypto.randomUUID(),
      fundId: params['fundId'] as string,
      fromOwner: body.fromOwner,
      toOwner: body.toOwner,
      units: body.units,
      transferredAt: new Date().toISOString(),
    };
    return HttpResponse.json(newTransfer, { status: 201 });
  }),
];
