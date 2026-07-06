"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useRouter } from "next/navigation";
import { api, getToken, setToken } from "@/lib/api";
import type {
  LoginRequest,
  LoginResponse,
  ProfileResponse,
  RegisterRequest,
  Role,
  SelectRoleRequest,
} from "@/lib/types";

interface AuthContextValue {
  profile: ProfileResponse | null;
  loading: boolean;
  isAuthenticated: boolean;
  needRoleSelection: boolean;
  availableRoles: Role[];
  login: (data: LoginRequest) => Promise<LoginResponse>;
  register: (data: RegisterRequest) => Promise<void>;
  selectRole: (role: Role) => Promise<void>;
  logout: () => Promise<void>;
  refreshProfile: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [profile, setProfile] = useState<ProfileResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [needRoleSelection, setNeedRoleSelection] = useState(false);
  const [availableRoles, setAvailableRoles] = useState<Role[]>([]);
  const router = useRouter();

  const refreshProfile = useCallback(async () => {
    const token = getToken();
    if (!token) {
      setProfile(null);
      setLoading(false);
      return;
    }
    try {
      const p = await api.get<ProfileResponse>("/profile");
      setProfile(p);
      setNeedRoleSelection(false);
    } catch {
      setToken(null);
      setProfile(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refreshProfile();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const login = useCallback(async (data: LoginRequest) => {
    const res = await api.post<LoginResponse>("/login", data, false);
    setToken(res.token);
    if (res.need_role_selection) {
      setNeedRoleSelection(true);
      setAvailableRoles(res.roles);
    } else {
      await refreshProfile();
    }
    return res;
  }, [refreshProfile]);

  const register = useCallback(async (data: RegisterRequest) => {
    await api.post("/register", data, false);
  }, []);

  const selectRole = useCallback(async (role: Role) => {
    const body: SelectRoleRequest = { role };
    const res = await api.post<LoginResponse>("/select-role", body);
    setToken(res.token);
    setNeedRoleSelection(false);
    await refreshProfile();
  }, [refreshProfile]);

  const logout = useCallback(async () => {
    try {
      await api.post("/logout");
    } catch {
      // ignore
    }
    setToken(null);
    setProfile(null);
    setNeedRoleSelection(false);
    router.push("/login");
  }, [router]);

  const value = useMemo<AuthContextValue>(
    () => ({
      profile,
      loading,
      isAuthenticated: !!profile,
      needRoleSelection,
      availableRoles,
      login,
      register,
      selectRole,
      logout,
      refreshProfile,
    }),
    [profile, loading, needRoleSelection, availableRoles, login, register, selectRole, logout, refreshProfile]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
