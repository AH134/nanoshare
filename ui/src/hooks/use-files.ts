import {
  queryOptions,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { deleteFile, getFiles, uploadFile } from "#/services/file";
import type { LinkPayload } from "#/services/link";

export const fileQueryOptions = queryOptions({
  queryKey: ["files"],
  queryFn: getFiles,
});

export function useFiles() {
  const queryClient = useQueryClient();
  const { data: files, isLoading } = useQuery(fileQueryOptions);

  const uploadMutation = useMutation({
    mutationFn: ({
      file,
      linkOptions,
    }: {
      file: File;
      linkOptions: LinkPayload;
    }) => uploadFile(file, linkOptions),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fileQueryOptions.queryKey });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteFile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fileQueryOptions.queryKey });
    },
  });

  return {
    files,
    isLoading,
    uploadMutation,
    deleteMutation,
  };
}
