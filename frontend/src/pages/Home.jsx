import React from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './Home.module.css'; // Подключаем CSS модуль
import logo from './logo.png'; // Импортируем изображение

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
                    <div className={styles.subtitle}>Найдите идеальный университет для себя</div> {/* Слоган */}
                </div>
                <div className={styles.searchWrapper}>
                    <input
                        type="text"
                        placeholder="Поиск..."
                        className={styles.searchInput}
                    />
                </div>
            </div>

            {/* основной контентик йоу какой крутой  */}
            <div className={styles.mainContentWrapper}>
                <div className={styles.filterSection}>
                    <div className={styles.filterTitle}>Фильтры</div>
                    {/* место для фильтров в будущем */}
                </div>

                <div className={styles.mainContent}>
                    <h1>Добро пожаловать!</h1>
                    <p>Мы рады видеть вас в нашем приложении. Пожалуйста, выберите одну из опций ниже:</p>
                    <div className={styles.buttonsWrapper}>
                        <button
                            className={styles.button}
                            onClick={() => navigate('/ProfileForm')}>
                            Создать профиль
                        </button>
                        <button
                            className={styles.button}
                            onClick={() => navigate('/login')}>
                            Войти
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Home;