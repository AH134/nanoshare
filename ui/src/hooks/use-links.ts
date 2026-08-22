import {
  queryOptions,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { createLink, getAllFileLinks, type LinkPayload } from "#/services/link";

export function linkQueryOptions(fileID: number) {
  return queryOptions({
    queryKey: ["links", fileID],
    queryFn: () => getAllFileLinks(fileID),
  });
}

export function useLinks(fileID: number) {
  const queryClient = useQueryClient();
  const { data: links, isLoading } = useQuery(linkQueryOptions(fileID));

  const createMutation = useMutation({
    mutationFn: (options: LinkPayload) => createLink(fileID, options),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: linkQueryOptions(fileID).queryKey,
      });
    },
  });

  return {
    links,
    isLoading,
    createMutation,
  };
}
