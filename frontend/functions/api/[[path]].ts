import { createProxyRequest, type PathParam } from './proxy';

type PagesContext = {
  request: Request;
  env: Record<string, string | undefined>;
  params: {
    path?: PathParam;
  };
};

export async function onRequest(context: PagesContext): Promise<Response> {
  try {
    return await fetch(createProxyRequest(context.request, context.env, context.params.path));
  } catch (error) {
    console.error(error);
    return Response.json({ error: 'API gateway is not configured.' }, { status: 500 });
  }
}
