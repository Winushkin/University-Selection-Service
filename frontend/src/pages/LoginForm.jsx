import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import  "./RegistrationForm.css"
import styles from "./Home.module.css";
import logo from "./logo.png";
import {useAuth} from "../AuthProvider.jsx";


function LogInForm() {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const { setAccessToken, setRefreshToken, setExpiresAt } = useAuth();
    const navigate = useNavigate();

    const handleLogin = async (e) => {
        e.preventDefault();

        try {
            const response = await fetch('/api/user/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ login, password })
            })



            if (!response.ok) {
                const errorData = await response.json();
                setError(errorData.message || 'Ошибка входа');
                return;
            }

            const data = await response.json();


            const expiresAt = Date.now() + 1000 * data.expires_in;
            localStorage.setItem('accessToken', data.access);
            localStorage.setItem('refreshToken', data.refresh);
            localStorage.setItem('expiresAt', expiresAt.toString());


            setAccessToken(data.access);
            setRefreshToken(data.refresh);
            setExpiresAt(expiresAt);


            navigate('/MainPage');
        } catch (err) {
            console.error('Login error:', err);
            setError('Ошибка соединения с сервером');
        }
    };


    return (
        <div>
            <div className={styles.header}>
                <div className={styles.logo}>
                    <img src={logo} alt="Logo" className={styles.logoImage}/>
                </div>
                <div className={styles.titleWrapper}>
                    <div className={styles.title}>UniQuest</div>
                    <div className={styles.subtitle}>Найдите идеальный университет для себя</div>
                </div>
            </div>

            <div className="register-container">
                <h2>Вход</h2>
                <form onSubmit={handleLogin}>
                    <div>
                        <label>Логин:</label>
                        <input
                            type="text"
                            value={login}
                            onChange={(e) => setLogin(e.target.value)}
                            required
                        />
                    </div>
                    <div>
                        <label>Пароль:</label>
                        <input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                    </div>
                    {error && <p style={{color: 'red'}}>{error}</p>}
                    <button type="submit">Войти</button>
                </form>
            </div>
        </div>
    );
}

export default LogInForm;



