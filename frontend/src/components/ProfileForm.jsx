import { useState } from 'react';

function ProfileForm() {
    const [profile, setProfile] = useState({
        egeScores: '',
        gpa: '',
        olympiads: '',
        desiredSpecialty: '',
        fieldOfActivity: '',
        educationType: 'очное',
        universityLocation: '',
        country: '',
        financing: 'Бюджет',
        importanceFactors: {
            universityRating: false,
            tuitionCost: false,
            dormitoryAvailability: false,
            location: false,
            scholarshipProgram: false,
        },
    });

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        if (type === 'checkbox') {
            setProfile((prevProfile) => ({
                ...prevProfile,
                importanceFactors: {
                    ...prevProfile.importanceFactors,
                    [name]: checked,
                },
            }));
        } else {
            setProfile((prevProfile) => ({
                ...prevProfile,
                [name]: value,
            }));
        }
    };

    return (
        <form>
            <div>
                <label>Баллы ЕГЭ:</label>
                <input type="text" name="egeScores" value={profile.egeScores} onChange={handleChange} />
            </div>
            <div>
                <label>Средний балл аттестата:</label>
                <input type="text" name="gpa" value={profile.gpa} onChange={handleChange} />
            </div>
            <div>
                <label>Наличие олимпиад:</label>
                <input type="text" name="olympiads" value={profile.olympiads} onChange={handleChange} />
            </div>
            <div>
                <label>Желаемая специальность:</label>
                <input type="text" name="desiredSpecialty" value={profile.desiredSpecialty} onChange={handleChange} />
            </div>
            <div>
                <label>Область деятельности:</label>
                <input type="text" name="fieldOfActivity" value={profile.fieldOfActivity} onChange={handleChange} />
            </div>
            <div>
                <label>Тип обучения:</label>
                <select name="educationType" value={profile.educationType} onChange={handleChange}>
                    <option value="очное">Очное</option>
                    <option value="заочное">Заочное</option>
                    <option value="дистанционное">Дистанционное</option>
                </select>
            </div>
            <div>
                <label>Город:</label>
                <input type="text" name="universityLocation" value={profile.universityLocation} onChange={handleChange} />
            </div>
            <div>
                <label>Страна:</label>
                <input type="text" name="country" value={profile.country} onChange={handleChange} />
            </div>
            <div>
                <label>Источник финансирования:</label>
                <select name="financing" value={profile.financing} onChange={handleChange}>
                    <option value="Бюджет">Бюджет</option>
                    <option value="Контракт">Контракт</option>
                </select>
            </div>
            <div>
                <label>Важность различных факторов:</label>
                <div>
                    <label>
                        <input type="checkbox" name="universityRating" checked={profile.importanceFactors.universityRating} onChange={handleChange} />
                        Рейтинг университета
                    </label>
                </div>
                <div>
                    <label>
                        <input type="checkbox" name="tuitionCost" checked={profile.importanceFactors.tuitionCost} onChange={handleChange} />
                        Стоимость обучения
                    </label>
                </div>
                <div>
                    <label>
                        <input type="checkbox" name="dormitoryAvailability" checked={profile.importanceFactors.dormitoryAvailability} onChange={handleChange} />
                        Наличие общежития
                    </label>
                </div>
                <div>
                    <label>
                        <input type="checkbox" name="location" checked={profile.importanceFactors.location} onChange={handleChange} />
                        Расположение
                    </label>
                </div>
                <div>
                    <label>
                        <input type="checkbox" name="scholarshipProgram" checked={profile.importanceFactors.scholarshipProgram} onChange={handleChange} />
                        Наличие стипендий
                    </label>
                </div>
            </div>
            <button type="button">Создать профиль</button>
        </form>
    );
}

export default ProfileForm;