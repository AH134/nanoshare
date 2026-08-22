import {
  queryOptions,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { getFiles, uploadFile } from "#/services/file";
import { createLink, type LinkPayload } from "#/services/link";
import { linkQueryOptions } from "./use-links";

export const fileQueryOptions = queryOptions({
  queryKey: ["files"],
  queryFn: getFiles,
});

export function useFiles() {
  const queryClient = useQueryClient();
  const { data: files, isLoading } = useQuery(fileQueryOptions);

  const uploadMutation = useMutation({
    mutationFn: async ({
      file,
      linkOptions,
    }: {
      file: File;
      linkOptions: LinkPayload;
    }) => {
      const uploadedFile = await uploadFile(file);
      await createLink(uploadedFile.id, linkOptions);
      return uploadedFile;
    },
    onSuccess: (uploadedFile) => {
      queryClient.invalidateQueries({ queryKey: fileQueryOptions.queryKey });
      queryClient.invalidateQueries({
        queryKey: linkQueryOptions(uploadedFile.id).queryKey,
      });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async () => {},
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
