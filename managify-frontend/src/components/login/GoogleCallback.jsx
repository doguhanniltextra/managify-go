import { useEffect, useContext, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { api } from "../api/api";
import { GOOGLE_AUTH_CALLBACK } from "../../constants/urls";
import { AuthContext } from "../../content/AuthContent";
import { toast } from "react-hot-toast";

export default function GoogleCallback() {
  const [searchParams] = useSearchParams();
  const { setToken } = useContext(AuthContext);
  const navigate = useNavigate();
  const [error, setError] = useState(null);

  useEffect(() => {
    const code = searchParams.get("code");

    if (!code) {
      setError("Authorization code not found");
      toast.error("Google login failed");
      setTimeout(() => navigate("/login"), 2000);
      return;
    }

    const handleCallback = async () => {
      try {
        const response = await api.post(GOOGLE_AUTH_CALLBACK, { code });
        localStorage.setItem("token", response.data.token);
        setToken(response.data.token);
        toast.success(`Welcome, ${response.data.name}!`);
        navigate("/dashboard");
      } catch (err) {
        setError("Authentication failed");
        toast.error("Google authentication failed");
        setTimeout(() => navigate("/login"), 2000);
      }
    };

    handleCallback();
  }, [searchParams, navigate, setToken]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        {error ? (
          <div>
            <p className="text-red-500 text-lg">{error}</p>
            <p className="text-gray-400 mt-2">Redirecting to login...</p>
          </div>
        ) : (
          <div>
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Signing in with Google...</p>
          </div>
        )}
      </div>
    </div>
  );
}
