import React, {
    createContext,
    useState,
    useEffect,
    useRef,
    useCallback,
    useContext
} from 'react';

// 1) Создаём контекст
const AuthContext = createContext();

// 2) Провайдер даёт в контекст токены и методы их обновления
export const AuthProvider = ({ children }) => {
    const [accessToken, setAccessToken]     = useState(localStorage.getItem('accessToken'));
    const [refreshToken, setRefreshToken]   = useState(localStorage.getItem('refreshToken'));
    const [expiresAt, setExpiresAt]         = useState(() => {
        const v = localStorage.getItem('expiresAt');
        return v ? Number(v) : null;
    });

    // Для хранения таймаута
    const refreshTimeout = useRef(null);

    // Функция обновления токена
    const refreshAccessToken = useCallback(async () => {
        try {
            const res = await fetch('/api/user/refresh', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh: refreshToken }),
            });
            const text = await res.text();
            const data = text ? JSON.parse(text) : {};
            if (!res.ok) throw new Error(data.message || `Ошибка ${res.status}`);

            const newAccess  = data.access;
            const newRefresh = data.refresh || refreshToken;
            const newExpire  = Date.now() + 15 * 60 * 1000;

            // Сохраняем в localStorage
            localStorage.setItem('accessToken',  newAccess);
            localStorage.setItem('refreshToken', newRefresh);
            localStorage.setItem('expiresAt',    String(newExpire));

            // Обновляем состояние
            setAccessToken(newAccess);
            setRefreshToken(newRefresh);
            setExpiresAt(newExpire);
        } catch (err) {
            console.error('Ошибка при обновлении токена:', err);
            // Очищаем при неуспехе
            localStorage.removeItem('accessToken');
            localStorage.removeItem('refreshToken');
            localStorage.removeItem('expiresAt');
            setAccessToken(null);
            setRefreshToken(null);
            setExpiresAt(null);
        }
    }, [refreshToken]);

    // Планируем авто-обновление
    useEffect(() => {
        if (refreshTimeout.current) {
            clearTimeout(refreshTimeout.current);
        }
        if (!refreshToken || !expiresAt) return;

        const delay = expiresAt - Date.now() - 60 * 1000;
        if (delay <= 0) {
            refreshAccessToken();
        } else {
            refreshTimeout.current = setTimeout(refreshAccessToken, delay);
        }

        return () => {
            if (refreshTimeout.current) clearTimeout(refreshTimeout.current);
        };
    }, [expiresAt, refreshToken, refreshAccessToken]);

    // 3) Передаём ВСЕ значения и сеттеры в контекст
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
};

// Хук для получения контекста
export const useAuth = () => useContext(AuthContext);
