import type { Person } from '$lib/types';
import type { PageLoad } from './$types';

export const load = (async () => {
    let persons: Person[] | null = await fetch('/persons').then(res => res.json());
    if (persons === null) persons = [];
    return { persons };
}) satisfies PageLoad;
