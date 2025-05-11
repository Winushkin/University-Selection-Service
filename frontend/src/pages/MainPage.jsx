import React, { useState } from 'react';
import styles from './MainPage.module.css';
import ToggleSwitch from '../components/ToggleSwitch';
import logo from "./logo.png";
import { useNavigate } from "react-router-dom";

function MainPage() {
    const navigate = useNavigate();
    const [error, setError] = useState('');

    const [importanceFactors, setImportanceFactors] = useState({
        ratingToPrestige: 1,
        ratingToEducationQuality: 1,
        ratingToScholarshipPrograms: 1,
        prestigeToEducationQuality: 1,
        prestigeToScholarshipPrograms: 1,
        educationQualityToScholarshipPrograms: 1,
        dormitory: false,
        scientificLabs: false,
        sportsInfrastructure: false,
        educationCost: '10000000'
    });

    const handleSliderChange = (e) => {
        const { name, value } = e.target;
        setImportanceFactors(prev => ({
            ...prev,
            [name]: Number(value)
        }));
    };

    const handleToggleChange = (key, value) => {
        setImportanceFactors(prev => ({
            ...prev,
            [key]: value
        }));
    };

    const handleChange = (e) => {
        const { name, value } = e.target;

        if (/^\d*$/.test(value)) {
            setImportanceFactors(prev => ({
                ...prev,
                [name]: value
            }));
        }
    };

    const handleEditProfile = () => {
        navigate('/EditProfileForm');
    };

    const handleLogout = () => {
        localStorage.clear();
        navigate('/');
    };

    const handleMAIRequest = async () => {
        try {
            const accessToken = localStorage.getItem('accessToken');
            const response = await fetch('http://localhost:80/api/analytic/analyze', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${accessToken}`,
                },
                body: JSON.stringify({
                    "ratingToPrestige": 0,
                    "ratingToEducationQuality": importanceFactors.ratingToEducationQuality,
                    "ratingToScholarshipPrograms": importanceFactors.ratingToScholarshipPrograms,
                    "prestigeToEducationQuality": importanceFactors.prestigeToEducationQuality,
                    "prestigeToScholarshipPrograms": importanceFactors.prestigeToScholarshipPrograms,
                    "educationQualityToScholarshipPrograms": importanceFactors.educationQualityToScholarshipPrograms,
                    "dormitory": importanceFactors.dormitory,
                    "scientificLabs": importanceFactors.scientificLabs,
                    "sportsInfrastructure": importanceFactors.sportsInfrastructure,
                    "educationCost": importanceFactors.educationCost
                }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                setError(errorData.message || 'Ошибка при запросе');
                return;
            }

            const data = await response.json();
            console.log('Успех:', data);

            const expiresAt = Date.now() + 1000 * 60 * 15;
            localStorage.setItem('accessToken', data.access);
            localStorage.setItem('refreshToken', data.refresh);
            localStorage.setItem('expiresAt', expiresAt);

            navigate('/MainPage');
        } catch (err) {
            setError('Ошибка соединения с сервером');
            console.error('Request error:', err);
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
                <button type="button" className="button" onClick={handleEditProfile}>
                    Редактировать профиль
                </button>
                <button type="button" className="button" onClick={handleLogout}>
                    Выйти
                </button>
            </div>

            <div className={styles.wrapper}>
                <aside className={styles.sidebar}>
                    <h2>Настройка фильтров</h2>

                    {/* Слайдеры */}
                    {[
                        { label: 'Рейтинг vs Престиж', name: 'ratingToPrestige' },
                        { label: 'Рейтинг vs Качество образования', name: 'ratingToEducationQuality' },
                        { label: 'Рейтинг vs Размер стипендии', name: 'ratingToScholarshipPrograms' },
                        { label: 'Престиж vs Качество образования', name: 'prestigeToEducationQuality' },
                        { label: 'Престиж vs Размер стипендии', name: 'prestigeToScholarshipPrograms' },
                        { label: 'Качество образования vs Размер стипендии', name: 'educationQualityToScholarshipPrograms' },
                    ].map(slider => (
                        <div className={styles.sliderGroup} key={slider.name}>
                            <label>
                                {slider.label}: <strong>{importanceFactors[slider.name]}</strong>
                            </label>
                            <input
                                type="range"
                                min="1"
                                max="9"
                                name={slider.name}
                                value={importanceFactors[slider.name]}
                                onChange={handleSliderChange}
                                className={styles.slider}
                            />
                        </div>
                    ))}

                    {/* Переключатели */}
                    <div className={styles.toggleGroup}>
                        <label>Общежитие</label>
                        <ToggleSwitch
                            checked={importanceFactors.dormitory}
                            onChange={val => handleToggleChange('dormitory', val)}
                        />
                    </div>
                    <div className={styles.toggleGroup}>
                        <label>Научные лаборатории</label>
                        <ToggleSwitch
                            checked={importanceFactors.scientificLabs}
                            onChange={val => handleToggleChange('scientificLabs', val)}
                        />
                    </div>
                    <div className={styles.toggleGroup}>
                        <label>Спортивная инфраструктура</label>
                        <ToggleSwitch
                            checked={importanceFactors.sportsInfrastructure}
                            onChange={val => handleToggleChange('sportsInfrastructure', val)}
                        />
                    </div>

                    <div className="form-group">
                        <label>Стоимость обучения должна быть не больше чем:</label>
                        <input
                            type="text"
                            name="educationCost"
                            className="form-input"
                            value={importanceFactors.educationCost}
                            onChange={handleChange}
                        />
                    </div>

                    {error && <div className="error-message">{error}</div>}

                    <button type="button" className="button" onClick={handleMAIRequest}>
                        Подобрать университет
                    </button>
                </aside>

                <main className={styles.mainContent}>
                    <h1>Список университетов</h1>
                    {/* Здесь будет отображаться результат фильтрации */}
                </main>
            </div>
        </div>
    );
}

export default MainPage;
