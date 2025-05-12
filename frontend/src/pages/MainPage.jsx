// src/pages/MainPage.jsx

import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../AuthProvider.jsx';
import ToggleSwitch from '../components/ToggleSwitch';
import styles from './MainPage.module.css';
import logo from './logo.png';

export default function MainPage() {
    const navigate = useNavigate();
    const { accessToken } = useAuth();
    const [error, setError] = useState('');
    const [universities, setUniversities] = useState([]);
    const [speciality, setSpeciality] = useState('');

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
        educationCost: '500000'
    });

    const handleSliderChange = e => {
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

    const handleChange = e => {
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
        setError('');
        try {
            if (!accessToken) {
                setError('Сначала нужно войти в систему');
                return;
            }

            const res = await fetch('/api/analytic/analyze', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${accessToken}`
                },
                body: JSON.stringify(importanceFactors)
            });

            if (!res.ok) {
                const err = await res.json().catch(() => ({}));
                setError(err.message || `Ошибка ${res.status}`);
                return;
            }

            const data = await res.json();
            console.log('Response data →', data)
            setSpeciality(data.speciality);
            setUniversities(data.universities);
        } catch (e) {
            console.error(e);
            setError('Ошибка соединения с сервером');
        }
    };

    return (
        <div>
            <header className={styles.header}>
                <div className={styles.logo}>
                    <img src={logo} alt="Logo" className={styles.logoImage} />
                </div>
                <div className={styles.titleWrapper}>
                    <h1 className={styles.title}>UniQuest</h1>
                    <p className={styles.subtitle}>Найдите идеальный университет для себя</p>
                </div>
                <button className="button" onClick={handleEditProfile}>
                    Редактировать профиль
                </button>
                <button className="button" onClick={handleLogout}>
                    Выйти
                </button>
            </header>

            <div className={styles.wrapper}>
                <aside className={styles.sidebar}>
                    <h2>Настройка фильтров</h2>

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
                        <label>Стоимость обучения не больше чем:</label>
                        <input
                            type="text"
                            name="educationCost"
                            className="form-input"
                            value={importanceFactors.educationCost}
                            onChange={handleChange}
                        />
                    </div>

                    {error && <div className="error-message">{error}</div>}

                    <button className="button" onClick={handleMAIRequest}>
                        Подобрать университет
                    </button>
                </aside>

                <main className={styles.mainContent}>
                    <h2>Результаты для специальности: {speciality || '—'}</h2>
                    {universities.length > 0 ? (
                        <ul className={styles.universityList}>
                            {universities.map((u, idx) => (
                                <li key={idx} className={styles.universityCard}>
                                    <h3>{u.name}</h3>
                                    <p>Регион: {u.region}</p>
                                    <p>Рейтинг: {u.rank}</p>
                                    <p>Стоимость: {u.cost}</p>
                                    <p>Престиж: {u.prestige}</p>
                                    <p>Качество образования: {u.quality}</p>
                                    <p>Общежитие: {u.dormitory ? 'да' : 'нет'}</p>
                                    <p>Лаборатории: {u.labs ? 'да' : 'нет'}</p>
                                    <p>Спорт. инфра-ры: {u.sport ? 'да' : 'нет'}</p>
                                    <p>Стипендия: {u.scholarship}</p>
                                    <p>Баллы ЕГЭ для бюджета: {u.BudgetPoints * 3}</p>
                                    <p>Баллы ЕГЭ для контракта: {u.ContractPoints * 3}</p>
                                    <p>Балл актуальности: {u.relevancy}</p>
                                    <p>
                                        Сайт:{' '}
                                        <a
                                            href={u.site}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className={styles.universityLink}
                                        >
                                            {u.site}
                                        </a>
                                    </p>
                                </li>
                            ))}
                        </ul>
                    ) : (
                        <p>Список университетов будет здесь после запроса.</p>
                    )}
                </main>
            </div>
        </div>
    );
}
