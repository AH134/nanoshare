import dayjs from "dayjs";
import type { LinkPayload } from "#/services/link";
import { fromDatetimeLocal, toDatetimeLocal } from "#/utils/date";

export interface LinkOptionsProps {
  options: LinkPayload;
  onChange: (options: LinkPayload) => void;
}

export function LinkOptions({ options, onChange }: LinkOptionsProps) {
  return (
    <div className="card bg-base-100 w-full border border-base-300">
      <div className="card-body">
        <div className="grid gap-2 sm:grid-cols-2 mb-1">
          <fieldset className="fieldset">
            <label className="label text-base-content" htmlFor="max-downloads">
              Max downloads
            </label>
            <input
              type="number"
              id="max-downloads"
              className="input"
              min={1}
              step={1}
              placeholder="Unlimited"
              value={options.maxDownloads ?? ""}
              onChange={(e) => {
                const val = e.target.value;
                onChange({
                  ...options,
                  maxDownloads: val === "" ? null : Number(val),
                });
              }}
            />
          </fieldset>
          <fieldset className="fieldset">
            <label className="label text-base-content" htmlFor="expiry-date">
              Expiry date test={new Date().toISOString()}---
              {/* {new Date().toUTCString()} */}
            </label>
            <input
              type="datetime-local"
              id="expiry-date"
              className="input"
              min={dayjs().format("YYYY-MM-DDTHH:mm")}
              value={toDatetimeLocal(options.expiresAt)}
              onChange={(e) =>
                onChange({
                  ...options,
                  expiresAt: fromDatetimeLocal(e.target.value),
                })
              }
            />
          </fieldset>
        </div>
        <p className="text-base-content/60 text-xs">
          These limits apply to files you add next. Leave max downloads or
          expiry empty for no limit.
        </p>
      </div>
    </div>
  );
}
