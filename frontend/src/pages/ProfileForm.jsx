import { useState } from 'react';
import React from 'react';
import { useNavigate } from 'react-router-dom';
import ToggleSwitch from '../components/ToggleSwitch.jsx';
import './ProfileForm.css';

function ProfileForm() {

    const Navigate=useNavigate();

    const [profile, setProfile] = useState({
        egeScores: '',
        gpa: '',

        desiredSpecialty: '',
        educationType: 'очное',
        country: '',
        universityLocation: '',
        financing: 'Бюджет',
        importanceFactors: {
            localUniversityRating: 5,
            prestige: 5,
            scholarshipPrograms: 5,
            educationQuality: 5,
            dormitory: false,
            scientificLabs: false,
            sportsInfrastructure: false
        },
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
            const response = await fetch('http://localhost:8080/profile', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(profile),
            });

            if (!response.ok) {
                throw new Error(`Ошибка сервера: ${response.status}`);
            }

            const data = await response.json();
            setSubmitStatus({ loading: false, success: true, error: null });
            console.log('Профиль успешно создан:', data);
            Navigate('/MainPage');
        } catch (error) {
            setSubmitStatus({ loading: false, success: false, error: error.message });
            console.error('Ошибка при отправке данных:', error);
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
                            <label>Средний балл аттестата:</label>
                            <input
                                type="text"
                                name="gpa"
                                className="form-input"
                                value={profile.gpa}
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
                            <label>Страна:</label>
                            <input
                                type="text"
                                name="country"
                                className="form-input"
                                value={profile.country}
                                onChange={handleChange}
                            />
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