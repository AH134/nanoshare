import {
	queryOptions,
	useMutation,
	useQuery,
	useQueryClient,
} from "@tanstack/react-query";
import { authService } from "#/services/auth";

export const authQueryOptions = queryOptions({
	queryKey: ["auth", "me"],
	queryFn: authService.getMe,
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
		}) => authService.login(username, password),
		onSuccess: (user) => {
			queryClient.setQueryData(authQueryOptions.queryKey, user);
		},
	});

	const logoutMutation = useMutation({
		mutationFn: authService.logout,
		onSettled: () => {
			queryClient.setQueryData(authQueryOptions.queryKey, null);
		},
	});

	return {
		user,
		isAuthenticated: !!user,
		isLoading,
		loginMutation,
		logoutMutation,
	};
}
