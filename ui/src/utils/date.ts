import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import updateLocale from "dayjs/plugin/updateLocale";

dayjs.extend(relativeTime);
dayjs.extend(updateLocale);
dayjs.updateLocale("en", {
  relativeTime: {
    future: "in %s",
    past: "%s ago",
    s: "just now",
    m: "1m",
    mm: "%dm",
    h: "1h",
    hh: "%dh",
    d: "1d",
    dd: "%dd",
    M: "1mo",
    MM: "%dmo",
    y: "1y",
    yy: "%dy",
  },
});

export function toDatetimeLocal(iso: string | null): string {
  return iso ? dayjs(iso).format("YYYY-MM-DDTHH:mm") : "";
}

export function fromDatetimeLocal(value: string): string | null {
  return value === "" ? null : dayjs(value).toISOString();
}

export default dayjs;
