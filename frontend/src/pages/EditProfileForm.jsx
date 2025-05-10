// src/pages/EditProfileForm.jsx

import React, { useState, useEffect, useRef } from 'react';
import './ProfileForm.css';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../AuthProvider.jsx';

export default function EditProfileForm() {
    const navigate = useNavigate();
    const { accessToken, refreshAccessToken } = useAuth();

    const didFetch = useRef(false);

    const [profile, setProfile] = useState({
        egeScores: '',
        desiredSpecialty: '',
        educationType: 'очное',
        universityLocation: '',
        financing: 'Бюджет',
    });

    const [loading, setLoading] = useState(true);

    const [submitStatus, setSubmitStatus] = useState({
        loading: false,
        success: false,
        error: null,
    });

    // Загрузка профиля один раз
    useEffect(() => {
        if (didFetch.current) return;
        didFetch.current = true;

        const init = async () => {
            try {
                // 1) Обновляем токен, если просрочен
                let token = accessToken;
                const expiresAt = Number(localStorage.getItem('expiresAt'));
                if (!token || Date.now() >= expiresAt) {
                    await refreshAccessToken();
                    token = localStorage.getItem('accessToken');
                }

                // 2) GET /api/user/profile
                const resp = await fetch('/api/user/profile', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${token}`,
                    },
                });

                if (!resp.ok) {
                    throw new Error(`Ошибка загрузки профиля: ${resp.status}`);
                }

                const data = await resp.json();

                // 3) Маппинг строго по ProfileResponse из proto
                setProfile({
                    egeScores: data.ege?.toString()      ?? '',
                    desiredSpecialty: data.speciality    ?? '',
                    educationType: data.eduType          ?? 'очное',
                    universityLocation: data.town        ?? '',
                    financing: data.financing            ?? 'Бюджет',
                });
            } catch (err) {
                setSubmitStatus(s => ({ ...s, error: err.message }));
            } finally {
                setLoading(false);
            }
        };

        init();
    }, [accessToken, refreshAccessToken]);

    // Обработчик изменения любого поля
    const handleChange = e => {
        const { name, value } = e.target;
        setProfile(prev => ({ ...prev, [name]: value }));
    };

    // Отправка обновлённого профиля
    const handleSubmit = async () => {
        setSubmitStatus({ loading: true, success: false, error: null });
        try {
            const token = localStorage.getItem('accessToken');
            if (!token) throw new Error('Нет accessToken для отправки');

            const resp = await fetch('/api/user/fill', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({
                    ege: Number(profile.egeScores),
                    speciality: profile.desiredSpecialty,
                    eduType: profile.educationType,
                    town: profile.universityLocation,
                    financing: profile.financing,
                }),
            });

            if (!resp.ok) {
                const errData = await resp.json();
                throw new Error(errData.message || `Ошибка обновления: ${resp.status}`);
            }

            await resp.json();
            setSubmitStatus({ loading: false, success: true, error: null });

            // Дадим пользователю секунду полюбоваться «успехом», затем уйдём
            setTimeout(() => {
                navigate('/MainPage');
            }, 1000);
        } catch (err) {
            setSubmitStatus({ loading: false, success: false, error: err.message });
        }
    };

    // Пока грузим
    if (loading) {
        return (
            <div className="form-container">
                <p>Загрузка профиля...</p>
                {submitStatus.error && (
                    <p className="error-message">{submitStatus.error}</p>
                )}
            </div>
        );
    }

    // Разметка формы
    return (
        <div className="form-container">
            <div className="form-card">
                <h2 className="form-title">Редактировать профиль</h2>
                <form onSubmit={e => e.preventDefault()}>
                    <div className="form-group">
                        <label>Баллы ЕГЭ:</label>
                        <input
                            type="text"
                            name="egeScores"
                            className="form-input"
                            value={profile.egeScores}
                            onChange={handleChange}
                        />
                    </div>

                    <div className="form-group">
                        <label>Специальность:</label>
                        <input
                            type="text"
                            name="desiredSpecialty"
                            className="form-input"
                            value={profile.desiredSpecialty}
                            onChange={handleChange}
                        />
                    </div>

                    <div className="form-group">
                        <label>Тип обучения:</label>
                        <select
                            name="educationType"
                            className="form-select"
                            value={profile.educationType}
                            onChange={handleChange}
                        >
                            <option value="очное">Очное</option>
                            <option value="заочное">Заочное</option>
                            <option value="дистанционное">Дистанционное</option>
                        </select>
                    </div>

                    <div className="form-group">
                        <label>Город:</label>
                        <input
                            type="text"
                            name="universityLocation"
                            className="form-input"
                            value={profile.universityLocation}
                            onChange={handleChange}
                        />
                    </div>

                    <div className="form-group">
                        <label>Источник финансирования:</label>
                        <select
                            name="financing"
                            className="form-select"
                            value={profile.financing}
                            onChange={handleChange}
                        >
                            <option value="Бюджет">Бюджет</option>
                            <option value="Контракт">Контракт</option>
                        </select>
                    </div>

                    {submitStatus.error && (
                        <div className="error-message" style={{ color: 'red' }}>
                            {submitStatus.error}
                        </div>
                    )}
                    {submitStatus.success && (
                        <div className="success-message" style={{ color: 'green' }}>
                            Изменения успешно сохранены!
                        </div>
                    )}

                    <button
                        type="button"
                        className="form-button"
                        onClick={handleSubmit}
                        disabled={submitStatus.loading}
                    >
                        {submitStatus.loading ? 'Сохранение...' : 'Сохранить изменения'}
                    </button>
                </form>
            </div>
        </div>
    );
}
