import { useParams } from "react-router";

export default function RoleReportsPage() {
  const { role } = useParams();

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          Reports - {role && role.charAt(0).toUpperCase() + role.slice(1)}
        </h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          View reports and analytics
        </p>
      </div>
      {/* Reports content will be added here */}
      <div className="bg-white rounded-xl p-6 dark:bg-white/[0.03]">
        <p className="text-gray-600 dark:text-gray-400">
          Reports content for {role} role
        </p>
      </div>
    </div>
  );
}
