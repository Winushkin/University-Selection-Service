import React from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './Home.module.css'; 
import logo from './logo.png'; 

//Функция стараницы регистрации и аутентификации
const Home = () => {
    const navigate = useNavigate();

    return (
        <div className={styles.pageWrapper}>
            {/* Header Section */}
            <div className={styles.header}>
                <div className={styles.logo}>
                    <img src={logo} alt="Logo" className={styles.logoImage} />
                </div>
                <div className={styles.titleWrapper}>
                    <div className={styles.title}>UniQuest</div>
                    <div className={styles.subtitle}>Найдите идеальный университет для себя</div>
                </div>

            </div>


            <div className={styles.mainContentWrapper}>


                <div className={styles.mainContent}>
                    <h1>Добро пожаловать!</h1>
                    <p>Мы рады видеть вас в нашем приложении. Пожалуйста, выберите одну из опций ниже:</p>
                    <div className={styles.buttonsWrapper}>
                        <button
                            className={styles.button}
                            onClick={() => navigate('/RegistrationForm')}>
                            Создать профиль
                        </button>
                        <button
                            className={styles.button}
                            onClick={() => navigate('/LoginForm')}>
                            Войти
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Home;
