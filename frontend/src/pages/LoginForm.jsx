import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import  "./RegistrationForm.css"
import styles from "./Home.module.css";
import logo from "./logo.png";


function LoginForm() {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleLogin = (e) => {
        e.preventDefault();

        /*
                const foundUser = users.find (
            (users) => users.login === login && users.password === password
        );

        if (!foundUser) {
            setError('Неверный логин или пароль');
            return;
        }

        console.log('Вход выполнен:', login);
        navigate('/'); //тут должна быть главная страница

            */
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

export default LoginForm;



