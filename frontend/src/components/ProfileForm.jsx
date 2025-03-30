import { useState } from 'react';
import ToggleSwitch from './ToggleSwitch';
import './ProfileForm.css';

function ProfileForm() {
    const [profile, setProfile] = useState({
        egeScores: '',
        gpa: '',
        olympiads: '',
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

    const handleToggleChange = (name, value) => {
        setProfile((prevProfile) => ({
            ...prevProfile,
            importanceFactors: {
                ...prevProfile.importanceFactors,
                [name]: value
            }
        }));
    };

    return (
        <div className="form-container">
            <div className="form-card">
                <h2 className="form-title">Профиль абитуриента</h2>
                <form>
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

                    <button type="button" className="form-button">Создать профиль</button>
                </form>
            </div>
        </div>
    );
}

export default ProfileForm;