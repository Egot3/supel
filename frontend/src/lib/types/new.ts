export type newRaw = {
	newId: string;
	userId: string;
	caption: string;
	bodyUrl: string;
	bodySize: number;
	imageUrls: string[];
	createdAt: string;
};

export type newCooked = {
	newId: string;
	userId: string;
	caption: string;
	body: Promise<string> | null;
	bodySize: number;
	imageUrls: string[];
	createdAt: string;
};
