type PagesMiddlewareContext = {
  request: Request;
  next: () => Promise<Response>;
};

const sensitiveJsonNames = new Set([
  'client_secret.json',
  'client_secrets.json',
  'credentials.json',
  'firebase-adminsdk.json',
  'firebase-credentials.json',
  'firebase-service-account.json',
  'firebase.json',
  'gcp-credentials.json',
  'gcp-service-account.json',
  'google-credentials.json',
  'google-service-account.json',
  'keyfile.json',
  'sa-key.json',
  'sa-private-key.json',
  'secrets.json',
  'service-account.json',
  'serviceaccountkey.json',
]);

export function isCredentialProbePath(pathname: string): boolean {
  const normalized = pathname.toLowerCase();
  const segments = normalized.split('/').filter(Boolean);
  const lastSegment = segments.length > 0 ? segments[segments.length - 1] : '';

  return (
    segments.some((segment) => segment === '.env' || segment.startsWith('.env.')) ||
    normalized.startsWith('/.aws/') ||
    sensitiveJsonNames.has(lastSegment)
  );
}

export async function onRequest(context: PagesMiddlewareContext): Promise<Response> {
  if (isCredentialProbePath(new URL(context.request.url).pathname)) {
    return new Response('Not found', {
      status: 404,
      headers: {
        'Content-Type': 'text/plain; charset=utf-8',
      },
    });
  }

  return context.next();
}
