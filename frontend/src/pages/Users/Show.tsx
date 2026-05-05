import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router";
import useGoBack from "../../hooks/useGoBack";

export default function UsersShow() {
  const { role, id } = useParams();
  const navigate = useNavigate();
  const goBack = useGoBack();
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // TODO: Fetch user data by id
    console.log("Fetching user with id:", id);
    setLoading(false);
  }, [id]);

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <p className="text-gray-500">Loading...</p>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="flex items-center justify-center p-8">
        <p className="text-red-500">User not found</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button
            onClick={goBack}
            className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
          >
            ← Back
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              User Details
            </h1>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              View user information
            </p>
          </div>
        </div>
        <button
          onClick={() => navigate(`/${role}/users/${id}/edit`)}
          className="px-4 py-2 bg-blue-500 rounded-lg font-medium text-white hover:bg-blue-600"
        >
          Edit
        </button>
      </div>

      {/* User Details */}
      <div className="bg-white rounded-xl p-6 dark:bg-white/[0.03]">
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
          {/* Username */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-400 mb-2">
              Username
            </label>
            <p className="px-4 py-2 bg-gray-50 rounded-lg dark:bg-gray-900 dark:text-white">
              {user.username}
            </p>
          </div>

          {/* Email */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-400 mb-2">
              Email
            </label>
            <p className="px-4 py-2 bg-gray-50 rounded-lg dark:bg-gray-900 dark:text-white">
              {user.email}
            </p>
          </div>

          {/* Phone */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-400 mb-2">
              Phone
            </label>
            <p className="px-4 py-2 bg-gray-50 rounded-lg dark:bg-gray-900 dark:text-white">
              {user.phone}
            </p>
          </div>

          {/* Role */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-400 mb-2">
              Role
            </label>
            <p className="px-4 py-2 bg-gray-50 rounded-lg dark:bg-gray-900 dark:text-white capitalize">
              {user.role?.replace(/_/g, " ")}
            </p>
          </div>

          {/* Status */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-400 mb-2">
              Status
            </label>
            <p className="px-4 py-2 bg-gray-50 rounded-lg dark:bg-gray-900 dark:text-white">
              {user.is_active ? "Active" : "Inactive"}
            </p>
          </div>

          {/* Created At */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-400 mb-2">
              Created At
            </label>
            <p className="px-4 py-2 bg-gray-50 rounded-lg dark:bg-gray-900 dark:text-white">
              {new Date(user.created_at).toLocaleString()}
            </p>
          </div>

          {/* Updated At */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-400 mb-2">
              Updated At
            </label>
            <p className="px-4 py-2 bg-gray-50 rounded-lg dark:bg-gray-900 dark:text-white">
              {new Date(user.updated_at).toLocaleString()}
            </p>
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-3 pt-8 border-t border-gray-200 dark:border-gray-700">
          <button
            onClick={goBack}
            className="px-6 py-2 border border-gray-300 rounded-lg font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-white/[0.05]"
          >
            Back
          </button>
        </div>
      </div>
    </div>
  );
}
