import axios from "axios";

type ApiErrorBody = {
  error?: {
    message?: string;
  };
  message?: string;
};

export function getApiErrorMessage(error: unknown, fallback: string) {
  if (!axios.isAxiosError<ApiErrorBody>(error)) {
    return fallback;
  }

  return error.response?.data?.error?.message ?? error.response?.data?.message ?? fallback;
}