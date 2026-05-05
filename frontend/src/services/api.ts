// API Client Base - Centralized API communication
import authService from "./auth";

const API_BASE_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1";

interface ApiResponse<T = unknown> {
  code?: number;
  message: string;
  data?: T;
  error?: unknown;
  errors?: Record<string, string>;
}

interface ApiError {
  status: number;
  message: string;
  error?: unknown;
  errors?: Record<string, string>;
}

export async function apiCall<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  const token = authService.getToken();

  let response: Response;
  let data: ApiResponse<T>;

  try {
    response = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options.headers,
      },
    });

    // Try to parse response as JSON
    try {
      data = await response.json();
    } catch (jsonError) {
      // If JSON parsing fails, it means backend returned non-JSON response
      const responseText = await response.text();
      console.error("Backend returned non-JSON response:", responseText);
      const error: ApiError = {
        status: response.status,
        message: `Server error: ${response.status} ${response.statusText}`,
        error: responseText,
        errors: undefined,
      };
      throw error;
    }
  } catch (err) {
    if (typeof err === "object" && err !== null && "status" in err) {
      throw err as ApiError;
    }
    // Network error or other fetch errors
    const error: ApiError = {
      status: 0,
      message: "Network error. Please check your connection.",
      error: err,
      errors: undefined,
    };
    throw error;
  }

  if (!response.ok) {
    const error: ApiError = {
      status: response.status,
      message: data.message || "An error occurred",
      error: data.error,
      errors: data.errors,
    };
    throw error;
  }

  return data.data as T;
}

export async function post<T>(endpoint: string, body: unknown): Promise<T> {
  return apiCall<T>(endpoint, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function get<T>(
  endpoint: string,
  headers?: Record<string, string>,
): Promise<T> {
  return apiCall<T>(endpoint, {
    method: "GET",
    headers: headers,
  });
}

export async function put<T>(endpoint: string, body: unknown): Promise<T> {
  return apiCall<T>(endpoint, {
    method: "PUT",
    body: JSON.stringify(body),
  });
}

export async function patch<T>(endpoint: string, body: unknown): Promise<T> {
  return apiCall<T>(endpoint, {
    method: "PATCH",
    body: JSON.stringify(body),
  });
}

export async function del<T>(endpoint: string): Promise<T> {
  return apiCall<T>(endpoint, {
    method: "DELETE",
  });
}
