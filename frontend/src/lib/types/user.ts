export type User = {
	uuid: string;
	nickname: string;
	avatarUrl: string;
	description: string | null;
	status: string | null;
	statusReactionKey: string | null;
};
