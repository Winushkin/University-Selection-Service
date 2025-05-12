// src/AuthProvider.jsx

import React, {
    createContext,
    useState,
    useEffect,
    useCallback,
    useContext
} from 'react';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
    // 1) Читаем токены и expiresAt из localStorage при старте
    const [accessToken, setAccessToken] = useState(
        () => localStorage.getItem('accessToken')
    );
    const [refreshToken, setRefreshToken] = useState(
        () => localStorage.getItem('refreshToken')
    );
    const [expiresAt, setExpiresAt] = useState(() => {
        const raw = localStorage.getItem('expiresAt');
        return raw ? Number(raw) : null;
    });

    // 2) Функция обновления
    const refreshAccessToken = useCallback(async () => {
        if (!refreshToken) return; // если нет refreshToken — нечего обновлять

        try {
            const res = await fetch('/api/user/refresh', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh: refreshToken })
            });

            if (res.status === 401) {
                // если сервер однозначно отверг refresh—очищаем токены
                throw new Error('refresh_token_invalid');
            }
            if (!res.ok) {
                throw new Error(`HTTP ${res.status}`);
            }

            const data = await res.json();
            const newAccess  = data.access;
            const newRefresh = data.refresh || refreshToken;
            const newExpire  = Date.now() + 10 * 60 * 1000; // +10 минут

            // сохраняем и в state, и в localStorage
            setAccessToken(newAccess);
            setRefreshToken(newRefresh);
            setExpiresAt(newExpire);

            localStorage.setItem('accessToken',  newAccess);
            localStorage.setItem('refreshToken', newRefresh);
            localStorage.setItem('expiresAt',    String(newExpire));
        } catch (err) {
            console.error('Не удалось обновить токен:', err);

            if (err.message === 'refresh_token_invalid') {
                // лишь при явном отказе сбрасываем
                setAccessToken(null);
                setRefreshToken(null);
                setExpiresAt(null);
                localStorage.removeItem('accessToken');
                localStorage.removeItem('refreshToken');
                localStorage.removeItem('expiresAt');
            }
            // при прочих ошибках (сеть, таймаут) — не трогаем токены
        }
    }, [refreshToken]);

    // 3) Запланируем авто-рефреш за минуту до expiresAt
    useEffect(() => {
        if (!expiresAt) return;
        const msUntilRefresh = expiresAt - Date.now() - 60 * 1000;
        if (msUntilRefresh <= 0) {
            // уже пора
            void refreshAccessToken();
            return;
        }
        const id = setTimeout(() => void refreshAccessToken(), msUntilRefresh);
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
