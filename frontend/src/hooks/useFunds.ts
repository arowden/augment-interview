import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

import { FundsService, type Fund, type CreateFundRequest } from '../api/client';

export const fundKeys = {
  all: ['funds'] as const,
  lists: () => [...fundKeys.all, 'list'] as const,
  list: () => [...fundKeys.lists()] as const,
  details: () => [...fundKeys.all, 'detail'] as const,
  detail: (id: string) => [...fundKeys.details(), id] as const,
};

export function useFunds() {
  return useQuery({
    queryKey: fundKeys.list(),
    queryFn: () => FundsService.listFunds(),
  });
}

export function useFund(id: string) {
  return useQuery({
    queryKey: fundKeys.detail(id),
    queryFn: () => FundsService.getFund(id),
    enabled: Boolean(id),
  });
}

export function useCreateFund() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateFundRequest) => FundsService.createFund(data),
    onSuccess: (newFund) => {
      queryClient.invalidateQueries({ queryKey: fundKeys.lists() });

      queryClient.setQueryData(fundKeys.detail(newFund.id), newFund);
    },
  });
}
