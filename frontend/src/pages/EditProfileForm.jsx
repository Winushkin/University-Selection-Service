
import React, { useState, useEffect, useRef } from 'react';
import './ProfileForm.css';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../AuthProvider.jsx';

// Списки регионов и специальностей
const REGIONS = [
    "Республика Тыва", "Московская область", "Смоленская область",
    "Орловская область", "Самарская область", "Республика Мордовия",
    "Севастополь", "Республика Калмыкия", "Карачаево-Черкесская Республика",
    "Чеченская Республика", "Амурская область", "Иркутская область",
    "Республика Хакасия", "Ивановская область", "Ханты-Мансийский автономный округ",
    "Вологодская область", "Новгородская область", "Курская область",
    "Санкт-Петербург", "Пермский край", "Саратовская область",
    "Республика Дагестан", "Свердловская область", "Рязанская область",
    "Брянская область", "Чувашская Республика", "Республика Крым",
    "Кабардино-Балкарская Республика", "Приморский край", "Тульская область",
    "Республика Марий Эл", "Кировская область", "Забайкальский край",
    "Магаданская область", "Москва", "Владимирская область",
    "Республика Коми", "Оренбургская область", "Республика Северная Осетия-Алания",
    "Томская область", "Кемеровская область", "Еврейская автономная область",
    "Республика Бурятия", "Белгородская область", "Калужская область",
    "Костромская область", "Калининградская область", "Нижегородская область",
    "Удмуртская Республика", "Волгоградская область", "Воронежская область",
    "Ярославская область", "Тверская область", "Пензенская область",
    "Ростовская область", "Республика Адыгея", "Тюменская область",
    "Тамбовская область", "Курганская область", "Камчатский край",
    "Архангельская область", "Мурманская область", "Республика Башкортостан",
    "Челябинская область", "Омская область", "Новосибирская область",
    "Алтайский край", "Республика Алтай", "Сахалинская область",
    "Липецкая область", "Астраханская область", "Краснодарский край",
    "Республика Саха (Якутия)", "Ленинградская область", "Псковская область",
    "Республика Карелия", "Республика Татарстан", "Ульяновская область",
    "Ставропольский край", "Республика Ингушетия", "Красноярский край",
    "Хабаровский край"
];

const SPECIALTIES = [
    "Юриспруденция", "Технологии легкой промышленности", "Геология",
    "Автоматика и управление", "Управление качеством",
    "Государственное и муниципальное управление", "Социальная работа",
    "Химическая и биотехнологии", "Управление водным транспортом",
    "Реклама и связи с общественностью", "Педагогическое образование",
    "Архитектура и градостроительство", "Авиационная и ракетно-космическая техника",
    "Энергетика и энергетическое машиностроение", "Транспортные средства",
    "Издательское дело", "Политология", "Металлургия",
    "Психолого-педагогическое и специальное (дефектологическое) образование",
    "Нефтегазовое дело", "Приборостроение и оптотехника",
    "Авиационные системы (эксплуатация)", "Пищевые технологии",
    "Электронная техника, радиотехника и связь", "Информационная безопасность",
    "Лесное дело", "Химия", "Социология", "Дизайн", "Математика",
    "Биология", "Менеджмент", "Здравоохранение", "Лингвистика и иностранные языки",
    "Экология", "Сельское и рыбное хозяйство", "Физическая культура",
    "Библиотеки и архивы", "Философия", "Строительство",
    "Геодезия и землеустройство", "Филология", "Психология",
    "Профессиональное обучение", "Морская техника",
    "Востоковедение и африканистика", "Машиностроение", "Материалы",
    "География", "Экономика", "Журналистика и литературное творчество",
    "История", "Информатика и вычислительная техника", "Сфера обслуживания",
    "Международные отношения", "Физика", "Технологические машины и оборудование",
    "Бизнес-информатика"
];


//Функция страницы редактирования профиля
export default function EditProfileForm() {
    const navigate = useNavigate();
    const { accessToken, refreshAccessToken } = useAuth();

    const didFetch = useRef(false);
    const [profile, setProfile] = useState({
        egeScores: '',
        desiredSpecialty: '',
        universityLocation: '',
        financing: 'Бюджет'
    });
    const [loading, setLoading] = useState(true);
    const [submitStatus, setSubmitStatus] = useState({
        loading: false,
        success: false,
        error: null
    });

    useEffect(() => {
        if (didFetch.current) return;
        didFetch.current = true;

        const init = async () => {
            try {
                // Обновляем токен, если просрочен
                let token = accessToken;


                // Запросим профиль
                const resp = await fetch('/api/user/profile', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${token}`
                    }
                });
                if (!resp.ok) {
                    throw new Error(`Ошибка загрузки профиля: ${resp.status}`);
                }
                const data = await resp.json();

                // Заполняем форму
                setProfile({
                    egeScores: data.ege?.toString()            ?? '',
                    desiredSpecialty: data.speciality          ?? '',
                    universityLocation: data.town              ?? '',
                    financing: data.financing                  ?? 'Бюджет'
                });
            } catch (err) {
                setSubmitStatus(s => ({ ...s, error: err.message }));
            } finally {
                setLoading(false);
            }
        };

        init();
    }, [accessToken, refreshAccessToken]);

    // Обработчик изменения полей
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
                    Authorization: `Bearer ${token}`
                },
                body: JSON.stringify({
                    ege: Number(profile.egeScores),
                    speciality: profile.desiredSpecialty,
                    town: profile.universityLocation,
                    financing: profile.financing
                })
            });
            if (!resp.ok) {
                const errData = await resp.json();
                throw new Error(errData.message || `Ошибка обновления: ${resp.status}`);
            }
            await resp.json();
            setSubmitStatus({ loading: false, success: true, error: null });

            // Покажем «успех» секунду, затем уйдём
            setTimeout(() => {
                navigate('/MainPage');
            }, 1000);
        } catch (err) {
            setSubmitStatus({ loading: false, success: false, error: err.message });
        }
    };

    if (loading) {
        return (
            <div className="form-container">
                <p>Загрузка профиля...</p>
                {submitStatus.error && <p className="error-message">{submitStatus.error}</p>}
            </div>
        );
    }

    return (
        <div className="form-container">
            <div className="form-card">
                <h2 className="form-title">Редактировать профиль</h2>
                <form onSubmit={e => e.preventDefault()}>
                    <div className="form-group">
                        <label>Среднее значение баллов по предметам ЕГЭ:</label>
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
                        <select
                            name="desiredSpecialty"
                            className="form-select"
                            value={profile.desiredSpecialty}
                            onChange={handleChange}
                        >
                            <option value="">— выберите специальность —</option>
                            {SPECIALTIES.map(spec => (
                                <option key={spec} value={spec}>{spec}</option>
                            ))}
                        </select>
                    </div>



                    <div className="form-group">
                        <label>Регион:</label>
                        <select
                            name="universityLocation"
                            className="form-select"
                            value={profile.universityLocation}
                            onChange={handleChange}
                        >
                            <option value="">— выберите регион —</option>
                            {REGIONS.map(region => (
                                <option key={region} value={region}>{region}</option>
                            ))}
                        </select>
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
                        <div className="error-message" style={{ color: 'red', marginBottom: '15px' }}>
                            {submitStatus.error}
                        </div>
                    )}
                    {submitStatus.success && (
                        <div className="success-message" style={{ color: 'green', marginBottom: '15px' }}>
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
