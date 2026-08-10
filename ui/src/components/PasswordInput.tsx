import { Eye, EyeOff } from "lucide-react";
import { useState } from "react";

interface PasswordInputProps {
  id: string;
  name: string;
  placeholder?: string;
  required?: boolean;
}

export default function PasswordInput({
  id,
  name,
  placeholder = "Enter your password",
  required,
}: PasswordInputProps) {
  const [showPassword, setShowPassword] = useState(false);

  const handleShowPassword = () => {
    setShowPassword((prev) => !prev);
  };
  return (
    <div className="relative">
      <input
        type={showPassword ? "text" : "password"}
        id={id}
        name={name}
        className="input bg-base-200/40 w-full pr-10 focus:outline-none"
        placeholder={placeholder}
        required={required}
      />
      <div
        className="absolute right-0 top-0 flex h-full items-center justify-center tooltip"
        data-tip={showPassword ? "Hide password" : "Show password"}
      >
        <button
          type="button"
          className="pr-2 h-full"
          tabIndex={-1}
          onClick={handleShowPassword}
        >
          {showPassword ? (
            <Eye className="size-4.5" />
          ) : (
            <EyeOff className="size-4.5" />
          )}
        </button>
      </div>
    </div>
  );
}
