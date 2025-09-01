import { useState, useEffect, useCallback } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { supabase } from '../services/supabase';

export function useAuth() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const location = useLocation();

  const isAuthenticated = useCallback(() => Boolean(user), [user]);

  const login = async (email, password) => {
    try {
      setLoading(true);
      setError(null);
      const { data, error: err } = await supabase.auth.signInWithPassword({ email, password });
      if (err) throw err;
      setUser(data.user || null);
      const from = location.state?.from?.pathname || '/';
      navigate(from, { replace: true });
      return { success: true };
    } catch (err) {
      const msg = err.message || 'Login failed';
      setError(msg);
      return { success: false, error: msg };
    } finally {
      setLoading(false);
    }
  };

  const logout = useCallback(async () => {
    await supabase.auth.signOut();
    setUser(null);
    navigate('/login');
  }, [navigate]);

  useEffect(() => {
    let mounted = true;
    const init = async () => {
      try {
        const { data } = await supabase.auth.getSession();
        if (mounted) setUser(data.session?.user || null);
      } finally {
        if (mounted) setLoading(false);
      }
    };
    init();
    const { data: sub } = supabase.auth.onAuthStateChange((_event, session) => {
      setUser(session?.user || null);
    });
    return () => { sub.subscription.unsubscribe(); mounted = false; };
  }, []);

  return { user, loading, error, isAuthenticated, login, logout };
}

export default useAuth;
