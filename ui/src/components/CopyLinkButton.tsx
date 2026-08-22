import { Check, Copy } from "lucide-react";
import { useState } from "react";

export interface CopyLinkButtonProps {
  url: string;
}

export function CopyLinkButton({ url }: CopyLinkButtonProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(url);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  return (
    <div className="tooltip" data-tip={copied ? "Copied!" : "Copy link"}>
      <button
        type="button"
        className={`btn btn-square btn-ghost ${copied && "bg-success/20"}`}
        onClick={handleCopy}
      >
        {copied ? (
          <Check className="size-4 text-success" />
        ) : (
          <Copy className="size-4" />
        )}
      </button>
    </div>
  );
}
