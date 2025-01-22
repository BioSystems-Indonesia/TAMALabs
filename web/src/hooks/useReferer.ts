import type { RedirectionSideEffect } from "react-admin";
import { useLocation } from "react-router-dom";

export function useRefererRedirect(defaultVal: RedirectionSideEffect) {
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const referer = queryParams.get("referer");

  return referer
    ? () => {
        return referer[0] === "/" ? referer.substring(1) : referer;
      }
    : defaultVal;
}

export function getRefererParam(): string {
  const location = useLocation();

  return generateParam(location.pathname);
}

function generateParam(location: string): string{
    return `referer=${location}`
}
