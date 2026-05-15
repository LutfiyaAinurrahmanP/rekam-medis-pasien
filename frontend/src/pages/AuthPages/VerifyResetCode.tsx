import PageMeta from "../../components/common/PageMeta";
import AuthLayout from "./AuthPageLayout";
import VerifyResetCodeForm from "../../components/auth/VerifyResetCodeForm";

export default function VerifyResetCode() {
  return (
    <>
      <PageMeta
        title="Verify Reset Code"
        description="Enter the reset code sent to your email to verify your identity."
      />
      <AuthLayout>
        <VerifyResetCodeForm />
      </AuthLayout>
    </>
  );
}
