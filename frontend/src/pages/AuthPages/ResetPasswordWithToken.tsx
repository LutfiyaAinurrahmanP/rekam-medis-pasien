import PageMeta from "../../components/common/PageMeta";
import AuthLayout from "./AuthPageLayout";
import ResetPasswordWithTokenForm from "../../components/auth/ResetPasswordWithTokenForm";

export default function ResetPasswordWithToken() {
  return (
    <>
      <PageMeta
        title="Reset Password"
        description="Set your new password using the reset token."
      />
      <AuthLayout>
        <ResetPasswordWithTokenForm />
      </AuthLayout>
    </>
  );
}
