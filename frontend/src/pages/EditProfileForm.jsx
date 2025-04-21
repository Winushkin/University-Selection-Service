import { useState, useEffect } from 'react';
import ToggleSwitch from '../components/ToggleSwitch.jsx';
import '../components/ProfileForm.css';

function EditProfileForm() {
    const [profile, setProfile] = useState(null);
    const [submitStatus, setSubmitStatus] = useState({
        loading: false,
        success: false,
        error: null
    });


    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const response = await fetch('http://localhost:8080/profile');
                if (!response.ok) throw new Error(`Ошибка загрузки: ${response.status}`);
                const data = await response.json();
                setProfile(data);
            } catch (error) {
                setSubmitStatus({ loading: false, success: false, error: error.message });
            }
        };

        fetchProfile();
    }, []);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;

        if (name.includes('.')) {
            const [parent, child] = name.split('.');
            setProfile((prev) => ({
                ...prev,
                [parent]: {
                    ...prev[parent],
                    [child]: type === 'checkbox' ? checked : Number(value),
                }
            }));
        } else {
            setProfile((prev) => ({
                ...prev,
                [name]: value
            }));
        }
    };

    const handleToggleChange = (name, value) => {
        setProfile((prev) => ({
            ...prev,
            importanceFactors: {
                ...prev.importanceFactors,
                [name]: value
            }
        }));
    };

    const handleSubmit = async () => {
        setSubmitStatus({ loading: true, success: false, error: null });

        try {
            const response = await fetch('http://localhost:8080/profile', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(profile),
            });

            if (!response.ok) throw new Error(`Ошибка обновления: ${response.status}`);
            const data = await response.json();
            setSubmitStatus({ loading: false, success: true, error: null });
            console.log('Профиль обновлён:', data);
        } catch (error) {
            setSubmitStatus({ loading: false, success: false, error: error.message });
            console.error('Ошибка при обновлении:', error);
        }
    };

    if (!profile) {
        return <div className="form-container">Загрузка профиля...</div>;
    }

    return (
        <div className="form-container">
            <div className="form-card">
                <h2 className="form-title">Редактировать профиль</h2>
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
                            <label>Наличие олимпиад:</label>
                            <input
                                type="text"
                                name="olympiads"
                                className="form-input"
                                value={profile.olympiads}
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

                    <div className="form-section">
                        <h3>Важность различных факторов</h3>
                        <div className="form-group">
                            <label>Местный рейтинг университета: {profile.importanceFactors.localUniversityRating}</label>
                            <input
                                type="range"
                                name="importanceFactors.localUniversityRating"
                                className="form-range"
                                min="1"
                                max="9"
                                value={profile.importanceFactors.localUniversityRating}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Престиж: {profile.importanceFactors.prestige}</label>
                            <input
                                type="range"
                                name="importanceFactors.prestige"
                                className="form-range"
                                min="1"
                                max="9"
                                value={profile.importanceFactors.prestige}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Стипендиальные программы: {profile.importanceFactors.scholarshipPrograms}</label>
                            <input
                                type="range"
                                name="importanceFactors.scholarshipPrograms"
                                className="form-range"
                                min="1"
                                max="9"
                                value={profile.importanceFactors.scholarshipPrograms}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Качество образования: {profile.importanceFactors.educationQuality}</label>
                            <input
                                type="range"
                                name="importanceFactors.educationQuality"
                                className="form-range"
                                min="1"
                                max="9"
                                value={profile.importanceFactors.educationQuality}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Наличие общежития:</label>
                            <ToggleSwitch
                                name="dormitory"
                                checked={profile.importanceFactors.dormitory}
                                onChange={(value) => handleToggleChange("dormitory", value)}
                            />
                        </div>
                        <div className="form-group">
                            <label>Наличие научных лабораторий:</label>
                            <ToggleSwitch
                                name="scientificLabs"
                                checked={profile.importanceFactors.scientificLabs}
                                onChange={(value) => handleToggleChange("scientificLabs", value)}
                            />
                        </div>
                        <div className="form-group">
                            <label>Наличие спортивной инфраструктуры:</label>
                            <ToggleSwitch
                                name="sportsInfrastructure"
                                checked={profile.importanceFactors.sportsInfrastructure}
                                onChange={(value) => handleToggleChange("sportsInfrastructure", value)}
                            />
                        </div>
                    </div>

                    {submitStatus.error && (
                        <div className="error-message" style={{ color: 'red' }}>
                            {submitStatus.error}
                        </div>
                    )}
                    {submitStatus.success && (
                        <div className="success-message" style={{ color: 'green' }}>
                            Профиль успешно обновлён!
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

export default EditProfileForm;
