import React, {
    createContext,
    useState,
    useEffect,
    useCallback,
    useContext
} from 'react';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
    const [accessToken, setAccessToken]   = useState(() => localStorage.getItem('accessToken'));
    const [refreshToken, setRefreshToken] = useState(() => localStorage.getItem('refreshToken'));
    const [expiresAt, setExpiresAt]       = useState(() => {
        const v = localStorage.getItem('expiresAt');
        return v ? Number(v) : null;
    });

    const refreshAccessToken = useCallback(async () => {
        if (!refreshToken) return;
        try {
            const res = await fetch('/api/user/refresh', {
                method: 'POST',
                headers: {'Content-Type':'application/json'},
                body: JSON.stringify({ refresh: refreshToken })
            });
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            const data = await res.json();
            const newAccess = data.access;
            const newRefresh = data.refresh  || refreshToken;
            const newExpires = Date.now() + 15 * 60 * 1000;

            setAccessToken(newAccess);
            setRefreshToken(newRefresh);
            setExpiresAt(newExpires);
            localStorage.setItem('accessToken',  newAccess);
            localStorage.setItem('refreshToken', newRefresh);
            localStorage.setItem('expiresAt',    String(newExpires));
        } catch (e) {
            console.error('Не удалось обновить токен:', e);
        }
    }, [refreshToken]);

    useEffect(() => {
        if (expiresAt && Date.now() >= expiresAt - 60 * 1000) {
            void refreshAccessToken();
        }
    }, [expiresAt, refreshAccessToken]);

   useEffect(() => {
        const intervalId = setInterval(() => {
            if (expiresAt && Date.now() >= expiresAt - 60 * 1000) {
                void refreshAccessToken();
            }
        }, 30 * 1000);

        return () => clearInterval(intervalId);
    }, [expiresAt, refreshAccessToken]);


    return (
        <AuthContext.Provider value={{
            accessToken,
            refreshToken,
            expiresAt,
            setAccessToken,
            setRefreshToken,
            setExpiresAt,
            refreshAccessToken
        }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    return useContext(AuthContext);
}
