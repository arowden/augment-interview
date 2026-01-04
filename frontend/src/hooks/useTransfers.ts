import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

import {
  TransfersService,
  type CreateTransferRequest,
  type CapTable,
  type CapTableEntry,
} from '../api/client';
import { capTableKeys } from './useCapTable';

export const transferKeys = {
  all: ['transfers'] as const,
  list: (fundId: string) => [...transferKeys.all, fundId] as const,
};

export function useTransfers(fundId: string) {
  return useQuery({
    queryKey: transferKeys.list(fundId),
    queryFn: () => TransfersService.listTransfers(fundId),
    enabled: Boolean(fundId),
  });
}

export function useCreateTransfer(fundId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTransferRequest) =>
      TransfersService.createTransfer(fundId, data),

    // Optimistic update.
    onMutate: async (newTransfer) => {
      // Cancel any outgoing refetches.
      await queryClient.cancelQueries({ queryKey: capTableKeys.detail(fundId) });
      await queryClient.cancelQueries({ queryKey: transferKeys.list(fundId) });

      // Snapshot the previous value.
      const previousCapTable = queryClient.getQueryData<CapTable>(
        capTableKeys.detail(fundId)
      );

      // Optimistically update the cap table.
      if (previousCapTable) {
        const updatedEntries = updateCapTableEntries(
          previousCapTable.entries,
          newTransfer.fromOwner,
          newTransfer.toOwner,
          newTransfer.units
        );

        queryClient.setQueryData<CapTable>(capTableKeys.detail(fundId), {
          ...previousCapTable,
          entries: updatedEntries,
        });
      }

      // Return context with previous value for rollback.
      return { previousCapTable };
    },

    // Rollback on error.
    onError: (_error, _newTransfer, context) => {
      if (context?.previousCapTable) {
        queryClient.setQueryData(
          capTableKeys.detail(fundId),
          context.previousCapTable
        );
      }
    },

    // Always refetch after mutation.
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: capTableKeys.detail(fundId) });
      queryClient.invalidateQueries({ queryKey: transferKeys.list(fundId) });
    },
  });
}

function updateCapTableEntries(
  entries: CapTableEntry[],
  fromOwner: string,
  toOwner: string,
  units: number
): CapTableEntry[] {
  const updatedEntries: CapTableEntry[] = [];
  let toOwnerExists = false;

  for (const entry of entries) {
    if (entry.ownerName === fromOwner) {
      // Deduct units from sender.
      const newUnits = entry.units - units;
      if (newUnits > 0) {
        updatedEntries.push({
          ...entry,
          units: newUnits,
          percentage: 0, // Will be recalculated.
        });
      }
      // If units become 0, don't include the entry.
    } else if (entry.ownerName === toOwner) {
      // Add units to recipient.
      toOwnerExists = true;
      updatedEntries.push({
        ...entry,
        units: entry.units + units,
        percentage: 0, // Will be recalculated.
      });
    } else {
      updatedEntries.push(entry);
    }
  }

  // If recipient doesn't exist, add new entry.
  if (!toOwnerExists) {
    updatedEntries.push({
      ownerName: toOwner,
      units,
      percentage: 0, // Will be recalculated.
      acquiredAt: new Date().toISOString(),
    });
  }

  // Recalculate percentages.
  const totalUnits = updatedEntries.reduce((sum, e) => sum + e.units, 0);
  return updatedEntries.map((entry) => ({
    ...entry,
    percentage: totalUnits > 0 ? (entry.units / totalUnits) * 100 : 0,
  }));
}
