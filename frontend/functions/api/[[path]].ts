import { createProxyRequest, shouldProxyApiPath, type PathParam } from './proxy';

type PagesContext = {
  request: Request;
  env: Record<string, string | undefined>;
  params: {
    path?: PathParam;
  };
};

export async function onRequest(context: PagesContext): Promise<Response> {
  if (!shouldProxyApiPath(context.params.path)) {
    return new Response('Not found', {
      status: 404,
      headers: {
        'Content-Type': 'text/plain; charset=utf-8',
      },
    });
  }

  try {
    return await fetch(createProxyRequest(context.request, context.env, context.params.path));
  } catch (error) {
    console.error(error);
    return Response.json({ error: 'API gateway is not configured.' }, { status: 500 });
  }
}
