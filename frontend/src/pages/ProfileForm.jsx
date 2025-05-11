import { useState } from 'react';
import React from 'react';
import { useNavigate } from 'react-router-dom';
import ToggleSwitch from '../components/ToggleSwitch.jsx';
import './ProfileForm.css';
import {useAuth} from "../AuthProvider.jsx";

function ProfileForm() {

    const Navigate=useNavigate();
    const [error, setError] = useState('');
    const { setAccessToken, setRefreshToken, setExpiresAt } = useAuth();

    const [profile, setProfile] = useState({
        egeScores: '',
        desiredSpecialty: '',
        educationType: 'очное',
        country: '',
        universityLocation: '',
        financing: 'Бюджет',

    });

    const [submitStatus, setSubmitStatus] = useState({
        loading: false,
        success: false,
        error: null
    });

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;

        if (name.includes('.')) {
            const [parent, child] = name.split('.');
            setProfile((prevProfile) => ({
                ...prevProfile,
                [parent]: {
                    ...prevProfile[parent],
                    [child]: type === 'checkbox' ? checked : Number(value),
                },
            }));
        } else {
            setProfile((prevProfile) => ({
                ...prevProfile,
                [name]: value,
            }));
        }
    };



    const handleSubmit = async () => {
        setSubmitStatus({ loading: true, success: false, error: null });


        try {
            const accessToken = localStorage.getItem('accessToken');
            const response = await fetch('http://localhost:80/api/user/fill', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${accessToken}`,
                },
                body: JSON.stringify({
                    ege: Number(profile.egeScores),
                    speciality: profile.desiredSpecialty,
                    eduType: profile.educationType,
                    town: profile.universityLocation,
                    financing: profile.financing,
                }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                setError(errorData.message || 'Ошибка при регистрации');
                return;
            }

            const data = await response.json();
            console.log('Регистрация успешна:', data);

            Navigate('/MainPage');
        } catch (err) {
            setError('Ошибка соединения с сервером');
            console.error('Signup error:', err);
        }
    };

    return (
        <div className="form-container">
            <div className="form-card">
                <h2 className="form-title">Профиль абитуриента</h2>
                <form onSubmit={(e) => e.preventDefault()}>
                    <div className="form-section">
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
                            <label>Cпециальность:</label>
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
                    </div>


                    {submitStatus.error && (
                        <div className="error-message" style={{color: 'red', marginBottom: '15px'}}>
                            {submitStatus.error}
                        </div>
                    )}

                    {submitStatus.success && (
                        <div className="success-message" style={{color: 'green', marginBottom: '15px'}}>
                            Профиль успешно создан!
                        </div>
                    )}

                    <button
                        type="button"
                        className="form-button"
                        onClick={handleSubmit}
                        disabled={submitStatus.loading}

                    >
                        {submitStatus.loading ? 'Отправка...' : 'Создать профиль'}

                    </button>


                </form>
            </div>
        </div>
    );
}

export default ProfileForm;