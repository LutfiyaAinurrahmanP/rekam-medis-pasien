// Auth Service - Business logic for authentication
import { post } from "./api";

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

class AuthService {
  async register(data: RegisterRequest): Promise<UserResponse> {
    return post<UserResponse>("/auth/register", data);
  }

  async login(data: LoginRequest): Promise<LoginResponse> {
    return post<LoginResponse>("/auth/login", data);
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
}

export default new AuthService();
