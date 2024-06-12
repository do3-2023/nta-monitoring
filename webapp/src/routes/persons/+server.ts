import { env } from '$env/dynamic/private';
import type { RequestHandler } from './$types';

export const POST = (async () => {
  let res = await fetch(`${env.API_URL ? env.API_URL : "http://localhost:3000"}/persons`, {
    method: 'POST'
  });
  console.log("New person status:", res.statusText);
  return res;
}) satisfies RequestHandler;

export const GET = (async () => {
  let res = await fetch(`${env.API_URL ? env.API_URL : "http://localhost:3000"}/persons`);
  console.log("Fetching persons:", res.statusText);
  return res;
}) satisfies RequestHandler;