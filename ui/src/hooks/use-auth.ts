import {
  queryOptions,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { ChangePassword, getMe, login, logout } from "#/services/auth";

export const authQueryOptions = queryOptions({
  queryKey: ["auth", "me"],
  queryFn: getMe,
  retry: false,
  staleTime: Infinity,
});

export function useAuth() {
  const queryClient = useQueryClient();
  const { data: user, isLoading } = useQuery(authQueryOptions);

  const loginMutation = useMutation({
    mutationFn: ({
      username,
      password,
    }: {
      username: string;
      password: string;
    }) => login(username, password),
    onSuccess: (user) => {
      queryClient.setQueryData(authQueryOptions.queryKey, user);
    },
  });

  const logoutMutation = useMutation({
    mutationFn: logout,
    onSettled: () => {
      queryClient.setQueryData(authQueryOptions.queryKey, null);
    },
  });

  const passwordMutation = useMutation({
    mutationFn: ChangePassword,
    onSuccess: () => {
      queryClient.clear();
    },
  });

  return {
    user,
    isAuthenticated: !!user,
    isLoading,
    loginMutation,
    logoutMutation,
    passwordMutation,
  };
}
