import { useQuery } from '@tanstack/react-query';

import { CapTableService } from '../api/client';

export const capTableKeys = {
  all: ['capTables'] as const,
  detail: (fundId: string) => [...capTableKeys.all, fundId] as const,
};

export function useCapTable(fundId: string) {
  return useQuery({
    queryKey: capTableKeys.detail(fundId),
    queryFn: () => CapTableService.getCapTable(fundId),
    enabled: Boolean(fundId),
  });
}
