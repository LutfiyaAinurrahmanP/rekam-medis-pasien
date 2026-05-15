import { useEffect, useState, useCallback } from "react";
import { get } from "../../../services/api";

import DetailModal from "../../../components/modals/DetailModal";
import { useLocation, useNavigate } from "react-router";
import { getRoleUsersPath } from "../../../pages/Roles/shared/role-routing";

interface ShowUserModalProps {
  isOpen: boolean;
  id: string | undefined;
  role?: string | undefined;
  onClose?: () => void; // optional callback
}

export default function ShowUserModal({
  isOpen,
  id,
  role,
  onClose,
}: ShowUserModalProps) {
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<Record<string, any> | null>(null);
  const navigate = useNavigate();
  const location = useLocation();

  const resolvedRole =
    role ??
    (location.pathname.split("/")[1]
      ? location.pathname.split("/")[1]
      : undefined);

  const close = useCallback(() => {
    if (onClose) {
      onClose();
      return;
    }

    // Fallback navigation only when no parent close handler is provided.
    navigate(getRoleUsersPath(resolvedRole));
  }, [navigate, onClose, resolvedRole]);

  useEffect(() => {
    let mounted = true;
    async function fetchUser() {
      if (!id) {
        setUser(null);
        setLoading(false);
        return;
      }
      setLoading(true);
      try {
        const data = await get<Record<string, any>>(`/users/${id}`);
        if (!mounted) return;
        setUser(data);
      } catch (err) {
        console.error("Failed to fetch user:", err);
        if (!mounted) return;
        setUser(null);
      } finally {
        if (!mounted) return;
        setLoading(false);
      }
    }

    if (isOpen) fetchUser();

    return () => {
      mounted = false;
    };
  }, [id, isOpen]);

  return (
    <DetailModal
      isOpen={isOpen}
      loading={loading}
      data={user}
      fields={[
        { key: "username", label: "Username" },
        { key: "email", label: "Email" },
        { key: "phone", label: "Phone" },
        {
          key: "role",
          label: "Role",
          render: (v) => String(v ?? "").replace(/_/g, " "),
        },
        {
          key: "is_active",
          label: "Status",
          render: (v) => (v ? "Active" : "Inactive"),
        },
        { key: "created_at", label: "Created At" },
        { key: "updated_at", label: "Updated At" },
      ]}
      onClose={close}
    />
  );
}
