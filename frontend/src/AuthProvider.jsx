// src/AuthProvider.jsx
import React, {
    createContext, useState, useEffect, useCallback, useContext
} from 'react';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
    // --- 1) Инициализация из localStorage
    const [accessToken, setAccessToken]   = useState(() => localStorage.getItem('accessToken'));
    const [refreshToken, setRefreshToken] = useState(() => localStorage.getItem('refreshToken'));
    const [expiresAt, setExpiresAt]       = useState(() => {
        const raw = localStorage.getItem('expiresAt');
        return raw ? Number(raw) : null;
    });

    // --- 2) Функция рефреша
    const refreshAccessToken = useCallback(async () => {
        if (!refreshToken) return;
        try {
            const res = await fetch('/api/user/refresh', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh: refreshToken })
            });
            if (!res.ok) throw new Error(`HTTP ${res.status}`);

            const data = await res.json();
            const newAccess  = data.access;
            const newRefresh = data.refresh  || refreshToken;
            const newExpire  = Date.now() + 15 * 60 * 1000; // +15 минут

            // --- 3) Сохраняем всё
            setAccessToken(newAccess);
            setRefreshToken(newRefresh);
            setExpiresAt(newExpire);

            localStorage.setItem('accessToken',  newAccess);
            localStorage.setItem('refreshToken', newRefresh);
            localStorage.setItem('expiresAt',    String(newExpire));
        } catch (err) {
            console.error('Не удалось обновить токен:', err);
            // здесь можно, по желанию, очищать токены при 401
        }
    }, [refreshToken]);

    // --- 4) Один единственный эффект с setTimeout
    useEffect(() => {
        if (!expiresAt) return;

        const msUntilRefresh = expiresAt - Date.now() - 60 * 1000;
        console.log('[Auth] expiresAt =', expiresAt, 'через мс до refresh =', msUntilRefresh);

        if (msUntilRefresh <= 0) {
            refreshAccessToken();
            return;
        }

        const id = setTimeout(() => {
            refreshAccessToken();
        }, msUntilRefresh);

        return () => clearTimeout(id);
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
