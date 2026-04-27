import { type User } from '$lib/types/user';
import { createContext } from 'svelte';

export const [getUserContext, setUserContext] = createContext<User>();
