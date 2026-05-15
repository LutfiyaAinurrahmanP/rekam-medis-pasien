// Auth Service - Business logic for authentication
import { post, get, put, patch, del } from "./api";

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

export interface UpdateMyProfileRequest {
  username?: string;
  email?: string;
  phone?: string;
  password?: string;
}

export interface ChangeMyPasswordRequest {
  old_password: string;
  new_password: string;
}

export interface DeleteMyAccountRequest {
  password: string;
  reason?: string;
}

export interface DeactivateMyAccountRequest {
  password: string;
  reason?: string;
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

  async getMe(): Promise<UserResponse> {
    return get<UserResponse>("/users/me");
  }

  async updateMe(data: UpdateMyProfileRequest): Promise<UserResponse> {
    return put<UserResponse>("/users/me", data);
  }

  async changeMyPassword(data: ChangeMyPasswordRequest): Promise<void> {
    await patch<void>("/users/me/change-password", data);
  }

  async deleteMe(data: DeleteMyAccountRequest): Promise<void> {
    await del<void>("/users/me", data);
  }

  async deactivateMe(data: DeactivateMyAccountRequest): Promise<void> {
    await patch<void>("/users/me/deactivate", data);
  }

  async getProfile(): Promise<UserResponse> {
    return this.getMe();
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
