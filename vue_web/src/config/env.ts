const normalizeOrigin = (value: string) => value.replace(/\/+$/, "");

const resolvedBackendOrigin =
  (import.meta.env.VITE_BACKEND_ORIGIN as string | undefined) ??
  (import.meta.env.DEV ? "http://localhost:8080" : window.location.origin);

export const backendOrigin = normalizeOrigin(resolvedBackendOrigin);
export const wsEndpoint = `${backendOrigin}/ws`;
