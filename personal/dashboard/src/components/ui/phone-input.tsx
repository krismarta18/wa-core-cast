"use client";

import { useState, useRef, useEffect } from "react";
import { ChevronDown } from "lucide-react";

export const COUNTRIES = [
  { code: "ID", name: "Indonesia", dial: "+62" },
  { code: "MY", name: "Malaysia", dial: "+60" },
  { code: "SG", name: "Singapura", dial: "+65" },
  { code: "PH", name: "Filipina", dial: "+63" },
  { code: "TH", name: "Thailand", dial: "+66" },
  { code: "VN", name: "Vietnam", dial: "+84" },
  { code: "MM", name: "Myanmar", dial: "+95" },
  { code: "CN", name: "China", dial: "+86" },
  { code: "JP", name: "Jepang", dial: "+81" },
  { code: "KR", name: "Korea Selatan", dial: "+82" },
  { code: "AU", name: "Australia", dial: "+61" },
  { code: "US", name: "United States", dial: "+1" },
  { code: "GB", name: "United Kingdom", dial: "+44" },
];

export function flagEmoji(code: string): string {
  return code
    .toUpperCase()
    .split("")
    .map((c) => String.fromCodePoint(0x1f1e6 + c.charCodeAt(0) - 65))
    .join("");
}

interface PhoneInputProps {
  value: string;
  onChange: (value: string) => void;
  dialCode: string;
  onDialCodeChange: (dialCode: string) => void;
  error?: string;
  id?: string;
}

export function PhoneInput({
  value,
  onChange,
  dialCode,
  onDialCodeChange,
  error,
  id,
}: PhoneInputProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const ref = useRef<HTMLDivElement>(null);

  const selected = COUNTRIES.find((c) => c.dial === dialCode) ?? COUNTRIES[0];

  const filtered = search.trim()
    ? COUNTRIES.filter(
        (c) =>
          c.name.toLowerCase().includes(search.toLowerCase()) ||
          c.dial.includes(search)
      )
    : COUNTRIES;

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
        setSearch("");
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  return (
    <div ref={ref} className="relative flex">
      {/* Dial code button */}
      <button
        type="button"
        onClick={() => {
          setOpen((o) => !o);
          setSearch("");
        }}
        className={`flex h-10 items-center gap-1.5 rounded-l-lg border border-r-0 bg-white px-3 text-sm text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-green-500 ${
          error ? "border-red-400" : "border-gray-300"
        }`}
      >
        <span className="text-base leading-none">{flagEmoji(selected.code)}</span>
        <span className="font-medium tabular-nums">{selected.dial}</span>
        <ChevronDown
          className={`h-3.5 w-3.5 text-gray-400 transition-transform duration-150 ${
            open ? "rotate-180" : ""
          }`}
        />
      </button>

      {/* Number input */}
      <input
        id={id}
        type="tel"
        inputMode="numeric"
        placeholder="8123456789"
        value={value}
        onChange={(e) => onChange(e.target.value.replace(/\D/g, ""))}
        className={`h-10 flex-1 rounded-r-lg border px-3 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500 ${
          error ? "border-red-400 focus:ring-red-400" : "border-gray-300"
        }`}
      />

      {/* Dropdown */}
      {open && (
        <div className="absolute left-0 top-full z-50 mt-1 w-72 overflow-hidden rounded-lg border border-gray-200 bg-white shadow-xl">
          {/* Search */}
          <div className="border-b border-gray-100 px-3 py-2">
            <input
              type="text"
              placeholder="Cari negara..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full rounded-md border border-gray-200 px-2 py-1.5 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-green-500"
              autoFocus
            />
          </div>
          <ul className="max-h-52 overflow-y-auto">
            {filtered.length === 0 ? (
              <li className="px-4 py-3 text-sm text-gray-400 text-center">
                Negara tidak ditemukan
              </li>
            ) : (
              filtered.map((c) => (
                <li key={c.code}>
                  <button
                    type="button"
                    onClick={() => {
                      onDialCodeChange(c.dial);
                      setOpen(false);
                      setSearch("");
                    }}
                    className={`flex w-full items-center gap-3 px-3 py-2 text-left text-sm transition-colors hover:bg-gray-50 ${
                      c.dial === dialCode
                        ? "bg-green-50 font-medium text-green-700"
                        : "text-gray-700"
                    }`}
                  >
                    <span className="text-lg leading-none">{flagEmoji(c.code)}</span>
                    <span className="flex-1 truncate">{c.name}</span>
                    <span className="tabular-nums text-gray-400">{c.dial}</span>
                  </button>
                </li>
              ))
            )}
          </ul>
        </div>
      )}
    </div>
  );
}
