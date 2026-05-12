import { useState, useEffect } from "react";
import { useParams } from "react-router";
import useGoBack from "../../../hooks/useGoBack";

export default function UsersEdit() {
  const { id } = useParams();
  const goBack = useGoBack();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    // TODO: Fetch user data by id
    console.log("Fetching user with id:", id);
  }, [id]);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    // TODO: Implement update user logic
    console.log("Update user with id:", id);
  };

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          onClick={goBack}
          className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
        >
          ← Back
        </button>
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Edit User
          </h1>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Update user information
          </p>
        </div>
      </div>

      {/* Form */}
      <div className="bg-white rounded-xl p-6 dark:bg-white/[0.03]">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            {/* Username */}
            <div>
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Username
              </label>
              <input
                type="text"
                className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
                placeholder="Enter username"
              />
            </div>

            {/* Email */}
            <div>
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Email
              </label>
              <input
                type="email"
                className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
                placeholder="Enter email"
              />
            </div>

            {/* Phone */}
            <div>
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Phone
              </label>
              <input
                type="tel"
                className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white"
                placeholder="Enter phone number"
              />
            </div>

            {/* Role */}
            <div>
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Role
              </label>
              <select className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white">
                <option>Select role</option>
                <option value="patient">Patient</option>
                <option value="doctor">Doctor</option>
                <option value="receptionist">Receptionist</option>
                <option value="admin">Admin</option>
                <option value="super_admin">Super Admin</option>
              </select>
            </div>

            {/* Is Active */}
            <div>
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Status
              </label>
              <select className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white">
                <option value="true">Active</option>
                <option value="false">Inactive</option>
              </select>
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={goBack}
              className="rounded-lg border border-gray-300 px-6 py-2 font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-white/[0.05]"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="rounded-lg bg-blue-500 px-6 py-2 font-medium text-white hover:bg-blue-600 disabled:opacity-50"
            >
              {loading ? "Updating..." : "Update User"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
