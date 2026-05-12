// Auth Service - Business logic for authentication
import { post, get } from "./api";

export interface RegisterRequest {
  username: string;
  email: string;
  phone: string;
  password: string;
  role?: string;
}

export interface UserResponse {
  id: number;
  username: string;
  email: string;
  phone: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username_or_email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  expires_at: string;
  user: UserResponse;
}

// Decode JWT to extract user info
function decodeJWT(token: string): Record<string, unknown> | null {
  try {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join(""),
    );
    return JSON.parse(jsonPayload);
  } catch (err) {
    console.error("Failed to decode JWT:", err);
    return null;
  }
}

class AuthService {
  async register(data: RegisterRequest): Promise<UserResponse> {
    return post<UserResponse>("/auth/register", data);
  }

  async login(data: LoginRequest): Promise<LoginResponse> {
    return post<LoginResponse>("/auth/login", data);
  }

  async getProfile(): Promise<UserResponse> {
    const token = this.getToken();
    if (!token) {
      throw new Error("No token found");
    }
    return get<UserResponse>("/users/me", {
      Authorization: `Bearer ${token}`,
    });
  }

  setToken(token: string): void {
    localStorage.setItem("authToken", token);
  }

  getToken(): string | null {
    return localStorage.getItem("authToken");
  }

  clearToken(): void {
    localStorage.removeItem("authToken");
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  // Helper to decode JWT payload
  getTokenPayload(): Record<string, unknown> | null {
    const token = this.getToken();
    if (!token) return null;
    return decodeJWT(token);
  }
}

export default new AuthService();
