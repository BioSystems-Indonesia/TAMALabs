import { useLocation } from "react-router-dom";

export default function useRefererRedirect(defaultVal: string) {
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const referer = queryParams.get("referer");

  return referer
    ? () => {
        return referer[0] === "/" ? referer.substring(1) : referer;
      }
    : defaultVal;
}
