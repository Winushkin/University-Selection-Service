import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import  "./RegistrationForm.css"
import styles from "./Home.module.css";
import logo from "./logo.png";

function RegistrationForm() {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleRegister = (e) => {
        e.preventDefault();

        if (password !== confirmPassword) {
            setError('Пароли не совпадают');
            return;
        }

        //  здесь отправить данные на сервер
        console.log('Регистрация успешна:', {login, password});

        navigate('/ProfileForm');
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

        <h2>Давай знакомиться!</h2>
        <form onSubmit={handleRegister}>
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
                <div>
                    <label>Повторите пароль:</label>
                    <input
                        type="password"
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        required
                    />
                </div>
                {error && <p style={{color: 'red'}}>{error}</p>}
                <button type="submit">Зарегистрироваться</button>
            </form>
        </div>
        </div>
    )
        ;
}

export default RegistrationForm;
