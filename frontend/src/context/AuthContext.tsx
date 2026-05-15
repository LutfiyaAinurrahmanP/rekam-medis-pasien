import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
  useCallback,
  useRef,
} from "react";
import authService from "../services/auth";

export interface User {
  id: number;
  username: string;
  email: string;
  phone: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

interface AuthContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  isAuthenticated: boolean;
  logout: () => void;
  loading: boolean;
  loadUserProfile: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [initialized, setInitialized] = useState(false);
  const profileFetchedRef = useRef(false);

  // Load user from localStorage on mount
  useEffect(() => {
    const storedUser = localStorage.getItem("currentUser");
    if (storedUser) {
      try {
        setUser(JSON.parse(storedUser));
      } catch (err) {
        console.error("Failed to parse stored user:", err);
        localStorage.removeItem("currentUser");
      }
    }
    setInitialized(true);
  }, []);

  // Save user to localStorage whenever it changes
  useEffect(() => {
    if (user) {
      localStorage.setItem("currentUser", JSON.stringify(user));
    } else {
      localStorage.removeItem("currentUser");
    }
  }, [user]);

  const loadUserProfile = useCallback(async () => {
    try {
      const profile = await authService.getMe();
      setUser(profile);
    } catch (err) {
      console.warn("Failed to load user profile from endpoint:", err);
      // Don't clear token/user on profile fetch failure, just log warning
      // Token might still be valid, user might already be set from login response
    }
  }, []);

  // Load user profile if token exists but user is not loaded (only once on init)
  useEffect(() => {
    if (initialized && !profileFetchedRef.current) {
      profileFetchedRef.current = true;

      if (authService.isAuthenticated() && !user) {
        loadUserProfile();
      }
      setLoading(false);
    }
  }, [initialized, loadUserProfile]);

  const logout = () => {
    setUser(null);
    authService.clearToken();
    localStorage.removeItem("currentUser");
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        setUser,
        isAuthenticated: !!user,
        logout,
        loading,
        loadUserProfile,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
