import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import "./RegistrationForm.css";
import styles from "./Home.module.css";
import logo from "./logo.png";
import { useAuth } from '../AuthProvider.jsx';

export default function RegistrationForm() {
    const navigate = useNavigate();


    const {
        setAccessToken,
        setRefreshToken,
        setExpiresAt
    } = useAuth();


    const [login, setLogin]                   = useState('');
    const [password, setPassword]             = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError]                   = useState('');


    const handleRegister = async (e) => {
        e.preventDefault();
        setError('');


        if (password !== confirmPassword) {
            setError('Пароли не совпадают');
            return;
        }

        try {

            const response = await fetch('/api/user/signup', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ login, password }),
            });


            const text = await response.text();
            let data;
            try {
                data = text ? JSON.parse(text) : {};
            } catch {
                throw new Error('Неправильный формат ответа сервера');
            }


            if (!response.ok) {
                setError(data.message || `Ошибка ${response.status}`);
                return;
            }




            const expiresAt = Date.now() + 15 * 60 * 1000;  // +15 минут
            localStorage.setItem('accessToken',  data.access);
            localStorage.setItem('refreshToken', data.refresh);
            localStorage.setItem('expiresAt',    String(expiresAt));


            setAccessToken(data.access);
            setRefreshToken(data.refresh);
            setExpiresAt(expiresAt);


            navigate('/ProfileForm');

        } catch (err) {
            console.error('Ошибка регистрации:', err);
            setError(err.message || 'Ошибка соединения с сервером');
        }
    };


    return (
        <div>
            <div className={styles.header}>
                <div className={styles.logo}>
                    <img src={logo} alt="Logo" className={styles.logoImage} />
                </div>
                <div className={styles.titleWrapper}>
                    <div className={styles.title}>UniQuest</div>
                    <div className={styles.subtitle}>Найдите идеальный университет для себя</div>
                </div>
            </div>

            <div className="register-container">
                <h2>Давай знакомиться!</h2>
                <form onSubmit={handleRegister}>
                    <div className="form-group">
                        <label>Логин:</label>
                        <input
                            type="text"
                            value={login}
                            onChange={e => setLogin(e.target.value)}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label>Пароль:</label>
                        <input
                            type="password"
                            value={password}
                            onChange={e => setPassword(e.target.value)}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label>Повторите пароль:</label>
                        <input
                            type="password"
                            value={confirmPassword}
                            onChange={e => setConfirmPassword(e.target.value)}
                            required
                        />
                    </div>

                    {error && <p className="error-message">{error}</p>}

                    <button type="submit" className="form-button">
                        Зарегистрироваться
                    </button>
                </form>
            </div>
        </div>
    );
}
